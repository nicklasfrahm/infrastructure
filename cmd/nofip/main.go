package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

// config defines the required configuration for this tool.
type config struct {
	// CloudflareAPIKey is the API key used to authenticate with Cloudflare.
	CloudflareAPIKey string
	// Record is the name of the DNS record to update.
	Record string
	// EdgeServers is the list of edge servers to use.
	EdgeServers []string
}

var cfg = &config{}

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

For this, it creates a DNS A records which points to
the IP addresses of the provided edge servers.

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
		fmt.Printf("Record name: %s\n", cfg.Record)
		fmt.Printf("Edge servers: %v\n", cfg.EdgeServers)

		if len(cfg.EdgeServers) == 0 {
			return errors.New("must specify at least one edge server")
		}

		// Use a waitgroup and a goroutine to resolve the IP addresses of
		// the edge servers and the current state of the edge record in
		// parallel.
		newV4EdgeIPs := make(map[string]bool)
		newV6EdgeIPs := make(map[string]bool)
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
						newV4EdgeIPs[ip.String()] = true
					} else {
						newV6EdgeIPs[ip.String()] = true
					}
				}
			}(server)
		}

		oldV4EdgeIPs := make(map[string]bool)
		oldV6EdgeIPs := make(map[string]bool)
		go func() {
			defer wg.Done()
			ips, err := net.LookupIP(cfg.Record)
			if err != nil {
				log.Printf("Failed to resolve edge record: %s: %s", cfg.Record, err)
				return
			}

			for _, ip := range ips {
				if ip.To4() != nil {
					oldV4EdgeIPs[ip.String()] = true
				} else {
					oldV6EdgeIPs[ip.String()] = true
				}
			}
		}()

		wg.Wait()

		fmt.Printf("[v4] old IPs: %+v\n", oldV4EdgeIPs)
		fmt.Printf("[v4] new IPs: %+v\n", newV4EdgeIPs)
		fmt.Printf("[v6] old IPs: %+v\n", oldV6EdgeIPs)
		fmt.Printf("[v6] new IPs: %+v\n", newV6EdgeIPs)

		if shouldUpdate(&oldV4EdgeIPs, &newV4EdgeIPs) {
			fmt.Println("[v4] should update")
			// TODO: Use the cloudflare SDK to create a DNS A record with the IPs.
		}

		if shouldUpdate(&oldV6EdgeIPs, &newV6EdgeIPs) {
			fmt.Println("[v6] should update")
			// TODO: Use the cloudflare SDK to create a DNS AAAA record with the IPs.
		}

		return nil
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

func shouldUpdate(old *map[string]bool, new *map[string]bool) bool {
	if len(*old) != len(*new) {
		return true
	}

	// Check if any new records should be created in the new set.
	for k := range *new {
		if !(*old)[k] {
			return true
		}
	}

	// Check if any old records should be removed in the new set.
	for k := range *old {
		if !(*new)[k] {
			return true
		}
	}

	return false
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
