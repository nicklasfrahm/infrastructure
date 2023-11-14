package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	cloudflare "github.com/cloudflare/cloudflare-go"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
)

// Config defines the required configuration for this tool.
type Config struct {
	// Record is the name of the DNS record to update.
	Record string
	// EdgeServers is the list of edge servers to use.
	EdgeServers []string
	// CloudflareToken is the Cloudflare API token to use.
	CloudflareToken string
}

var cfg = &Config{}

var version = "dev"
var help bool
var noWait bool

var rootCmd = &cobra.Command{
	Use:   "nofip",
	Short: "Expose a service without floating IPs",
	Long: `               __ _
  _ __   ___  / _(_)_ __
 | '_ \ / _ \| |_| | '_ \
 | | | | (_) |  _| | |_) |
 |_| |_|\___/|_| |_| .__/
                   |_|

nofip is a tool that allows you to expose a service
using high-availability without using floating IPs.

For this, it creates multiple DNS A records which
point to the IP addresses of the provided edge servers.

If you set the service edge to be europe.example.com
and that your edge servers are alfa.example.com and
bravo.example.com, it will create a DNS A record for
europe.example.com with the IP addresses of alfa and
bravo.

Note that this can't be used for services that require
geo-DNS, as nofip will only use simple A records.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if help {
			cmd.Help()
			os.Exit(0)
		}
	},
	RunE:         reconcileEdgeRecord,
	Version:      version,
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "display help for command")
	rootCmd.PersistentFlags().StringVarP(&cfg.Record, "record", "r", "", "name of the DNS record to update")
	rootCmd.MarkPersistentFlagRequired("record")
	rootCmd.PersistentFlags().StringSliceVarP(&cfg.EdgeServers, "edge-servers", "e", []string{}, "list of edge servers (comma-separated)")
	rootCmd.MarkPersistentFlagRequired("edge-servers")
	rootCmd.PersistentFlags().BoolVarP(&noWait, "no-wait", "n", false, "don't wait for DNS propagation")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func createRecord(api *cloudflare.API, zoneID string, ip string, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	fmt.Printf("[++] Creating: %s\n", ip)
	_, err := api.CreateDNSRecord(context.TODO(), cloudflare.ResourceIdentifier(zoneID), cloudflare.CreateDNSRecordParams{
		Name:    cfg.Record,
		Type:    "A",
		Content: ip,
		Proxied: cloudflare.BoolPtr(false),
		TTL:     60,
	})
	if err != nil {
		if e, ok := err.(*cloudflare.RequestError); ok {
			// Ignore conflicting record errors.
			if !e.InternalErrorCodeIs(81057) {
				fmt.Printf("[v4] Failed to create A record: %s", err)
			}
		}
	}
}

func deleteRecord(api *cloudflare.API, record *cloudflare.DNSRecord, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	fmt.Printf("[--] Deleting: %s\n", record.Content)
	err := api.DeleteDNSRecord(context.TODO(), cloudflare.ResourceIdentifier(record.ZoneID), record.ID)
	if err != nil {
		fmt.Printf("[v4] Failed to delete A record: %s", err)
	}
}

func purgeResolverCache() error {
	fmt.Printf("[**] Purging resolver cache at 1.1.1.1 ...\n")
	_, err := http.Post(fmt.Sprintf("https://1.1.1.1/api/v1/purge?type=A&domain=%s", cfg.Record), "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to purge resolve cache: %s", err)
	}

	return nil
}

func reconcileEdgeRecord(cmd *cobra.Command, args []string) error {
	// Verify that the Cloudflare API token is set.
	cfg.CloudflareToken = os.Getenv("CLOUDFLARE_TOKEN")
	if cfg.CloudflareToken == "" {
		return errors.New("missing environment variable: CLOUDFLARE_TOKEN")
	}

	// Ensure that the record follows the storage format of records in Cloudflare.
	if cfg.Record[len(cfg.Record)-1] == '.' {
		cfg.Record = cfg.Record[:len(cfg.Record)-1]
	}
	cfg.Record = strings.ToLower(cfg.Record)

	fmt.Printf("[>>] Record name: %s\n", cfg.Record)
	fmt.Printf("[>>] Edge servers: %v\n", cfg.EdgeServers)

	if len(cfg.EdgeServers) == 0 {
		return errors.New("must specify at least one edge server")
	}

	v4Reconciler := NewReconciler()
	v6Reconciler := NewReconciler()

	// Use a waitgroup and a goroutine to resolve the IP addresses of
	// the edge servers and the current state of the edge record in
	// parallel.
	wg := sync.WaitGroup{}
	wg.Add(len(cfg.EdgeServers) + 1)

	for _, server := range cfg.EdgeServers {
		go func(server string) {
			defer wg.Done()
			ips, err := net.LookupIP(server)
			if err != nil {
				log.Printf("Failed to resolve edge server: %s: %s", server, err)
				return
			}

			for _, ip := range ips {
				if ip.To4() != nil {
					v4Reconciler.Desired(ip.String())
				} else {
					v6Reconciler.Desired(ip.String())
				}
			}
		}(server)
	}

	// TODO: Use cloudflare API instead of DNS resolver to
	// avoid unnecessary reconciliation due to DNS caching.
	go func() {
		defer wg.Done()
		ips, err := net.LookupIP(cfg.Record)
		if err != nil {
			fmt.Printf("[XX] Failed to resolve edge record: %s: %s", cfg.Record, err)
			return
		}

		for _, ip := range ips {
			if ip.To4() != nil {
				v4Reconciler.Current(ip.String())
			} else {
				v6Reconciler.Current(ip.String())
			}
		}
	}()

	wg.Wait()

	api, err := cloudflare.NewWithAPIToken(cfg.CloudflareToken)
	if err != nil {
		return fmt.Errorf("failed to create Cloudflare API client: %s", err)
	}

	zoneRegex := regexp.MustCompile("[a-z]+.[a-z]+$")
	zoneName := zoneRegex.FindString(cfg.Record)
	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return fmt.Errorf("failed to list zones: %s", err)
	}

	if v4Reconciler.ShouldUpdate() {
		if len(v4Reconciler.Plan.Create) > 0 {
			fmt.Printf("[v4] Creating %d records ...\n", len(v4Reconciler.Plan.Create))

			wgCreate := &sync.WaitGroup{}
			for ip := range v4Reconciler.Plan.Create {
				wgCreate.Add(1)
				go createRecord(api, zoneID, ip, wgCreate)
			}
			wgCreate.Wait()
		}

		if len(v4Reconciler.Plan.Delete) > 0 {
			fmt.Printf("[v4] Deleting %d records ...\n", len(v4Reconciler.Plan.Delete))

			records, _, err := api.ListDNSRecords(context.TODO(), cloudflare.ResourceIdentifier(zoneID), cloudflare.ListDNSRecordsParams{
				Type: "A",
				Name: cfg.Record,
			})
			if err != nil {
				return fmt.Errorf("failed to list relevant records: %s", err)
			}

			wgDelete := &sync.WaitGroup{}
			for _, record := range records {
				// Check if the record should be deleted.
				if v4Reconciler.Plan.Delete[record.Content] {
					wgDelete.Add(1)
					go deleteRecord(api, &record, wgDelete)
				}
			}
			wgDelete.Wait()
		}
	}

	if v6Reconciler.ShouldUpdate() {
		return errors.New("IPv6 is not supported yet")
		// TODO: Use the cloudflare SDK to create a DNS AAAA record with the IPs.
	}

	if err := purgeResolverCache(); err != nil {
		return err
	}

	// TODO: Implement wait for DNS propagation.

	return nil
}
