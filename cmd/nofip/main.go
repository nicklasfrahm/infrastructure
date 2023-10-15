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

// ReconciliationState defines the state of a resource
// before the reconciliation process.
type ReconciliationState struct {
	// Current is the current state of the resource.
	Current map[string]bool
	// Desired is the desired state of the resource.
	Desired map[string]bool
}

// ReconciliationPlan defines the plan to reconcile a
// resource.
type ReconciliationPlan struct {
	// Create is the list of records to create.
	Create map[string]bool
	// Delete is the list of records to delete.
	Delete map[string]bool
}

// Reconciliation is stores the current and desired state
// of a resource along with a plan to reconcile them.
type Reconciliation struct {
	sync.Mutex
	// State is the state of the resource before the reconciliation.
	State *ReconciliationState
	// Plan is the plan to reconcile the resource.
	Plan *ReconciliationPlan
}

// ShouldUpdate updates the plan to reconcile the resource.
func (r *Reconciliation) ShouldUpdate() bool {
	shouldUpdate := false

	// Check if any new records should be created in the new set.
	for key := range r.State.Desired {
		if !r.State.Current[key] {
			shouldUpdate = true
			r.Plan.Create[key] = true
		}
	}

	// Check if any old records should be removed in the new set.
	for key := range r.State.Current {
		if !r.State.Desired[key] {
			shouldUpdate = true
			r.Plan.Delete[key] = true
		}
	}

	return shouldUpdate
}

// NewReconciliation creates a new reconciliation object.
func NewReconciliation() *Reconciliation {
	return &Reconciliation{
		State: &ReconciliationState{
			Current: map[string]bool{},
			Desired: map[string]bool{},
		},
		Plan: &ReconciliationPlan{
			Create: map[string]bool{},
			Delete: map[string]bool{},
		},
	}
}

var cfg = &Config{}

var version = "dev"
var help bool

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
	RunE: func(cmd *cobra.Command, args []string) error {
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

		v4Update := NewReconciliation()
		v6Update := NewReconciliation()

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
						v4Update.Lock()
						v4Update.State.Desired[ip.String()] = true
						v4Update.Unlock()
					} else {
						v6Update.Lock()
						v6Update.State.Desired[ip.String()] = true
						v6Update.Unlock()
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
				fmt.Printf("[] Failed to resolve edge record: %s: %s", cfg.Record, err)
				return
			}

			for _, ip := range ips {
				if ip.To4() != nil {
					v4Update.Lock()
					v4Update.State.Current[ip.String()] = true
					v4Update.Unlock()
				} else {
					v6Update.Lock()
					v6Update.State.Current[ip.String()] = true
					v6Update.Unlock()
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

		if v4Update.ShouldUpdate() {
			if len(v4Update.Plan.Create) > 0 {
				fmt.Printf("[v4] Creating %d records ...\n", len(v4Update.Plan.Create))

				wgCreate := &sync.WaitGroup{}
				for ip := range v4Update.Plan.Create {
					wgCreate.Add(1)
					go createRecord(api, zoneID, ip, wgCreate)
				}
				wgCreate.Wait()
			}

			if len(v4Update.Plan.Delete) > 0 {
				fmt.Printf("[v4] Deleting %d records ...\n", len(v4Update.Plan.Delete))

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
					if v4Update.Plan.Delete[record.Content] {
						wgDelete.Add(1)
						go deleteRecord(api, &record, wgDelete)
					}
				}
				wgDelete.Wait()
			}
		}

		if v6Update.ShouldUpdate() {
			fmt.Println("[v6] should update")
			// TODO: Use the cloudflare SDK to create a DNS AAAA record with the IPs.
		}

		return purgeResolverCache()
	},
	Version:      version,
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "display help for command")
	rootCmd.PersistentFlags().StringVarP(&cfg.Record, "record", "r", "", "name of the DNS record to update")
	rootCmd.MarkPersistentFlagRequired("record")
	rootCmd.PersistentFlags().StringSliceVarP(&cfg.EdgeServers, "edge-servers", "e", []string{}, "list of edge servers (comma-separated)")
	rootCmd.MarkPersistentFlagRequired("edge-servers")
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
