package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const (
	// EndpointIPv4 is an endpoint to fetch the public IPv4 address.
	EndpointIPv4 = "https://ip4.seeip.org"
	// EndpointIPv6 is an endpoint to fetch the public IPv6 address.
	EndpointIPv6 = "https://ip6.seeip.org"
)

// RecordType is the type of DNS record.
type RecordType string

const (
	// RecordTypeA is the A record type.
	RecordTypeA RecordType = "A"
	// RecordTypeAAAA is the AAAA record type.
	RecordTypeAAAA RecordType = "AAAA"
)

var ddnsCmd = &cobra.Command{
	Use:   "ddns <domain>",
	Short: "Dynamically update DNS records",
	Long: `A dynamic DNS client to update or create Google Cloud DNS records.

For this command to work, the service account used needs to have
the DNS Admin role. You may get away with fewer permissions but
this has not been tested.

The recommended usage of this command is as part of a cron job,
either using crontab or a Kubernetes CronJob.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Configure logger.
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		})

		// Parse command line arguments.
		domain := CanonicalDomain(args[0])

		// Fetch Google Cloud credentials and use them
		// to create a service for the Cloud DNS API.
		provider, err := NewGoogleProvider(domain)
		if err != nil {
			return err
		}

		// Get the current IP.
		ipv4, err := PublicIP(EndpointIPv4)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch IPv4 address")
			return err
		}
		ipv6, err := PublicIP(EndpointIPv6)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch IPv6 address")
			return err
		}

		// Reconcile the records.
		if err := provider.Reconcile(domain, RecordTypeA, ipv4); err != nil {
			return err
		}
		if err := provider.Reconcile(domain, RecordTypeAAAA, ipv6); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ddnsCmd)
}

// CanonicalDomain ensures that the domain has
// the canonical format with the trailing dot.
func CanonicalDomain(domain string) string {
	// Using strings.HasSuffix() seems overkill here.
	if domain[len(domain)-1] != '.' {
		return domain + "."
	}

	return domain
}

// PublicIP fetches the public IP address from
// the given endpoint using HTTP. Note that this
// function will NOT fail if there is no source
// address available.
func PublicIP(endpoint string) (net.IP, error) {
	// Fetch public IP address.
	resp, err := http.Get(endpoint)
	if err != nil {
		var opError *net.OpError
		if errors.As(err, &opError) {
			// Return no IP if we don't have
			// an address for this IP family.
			if opError.Source == nil {
				return nil, nil
			}
		}
		return nil, err
	}

	// Read response.
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse IP address.
	return net.ParseIP(string(ip)), nil
}

// Provider describes the interface of a DNS record provider.
type Provider interface {
	Reconcile(domain, recordType, desired string) error
}

// GoogleCredentials describes a subset of the
// structure of the Google Cloud credentials.
type GoogleCredentials struct {
	ProjectID string `json:"project_id"`
}

// GoogleProvider implements the Provider interface
// for Google Cloud DNS.
type GoogleProvider struct {
	ProjectID string
	ZoneID    string
	service   *dns.Service
}

// NewGoogleProvider creates a new GoogleProvider.
func NewGoogleProvider(domain string) (*GoogleProvider, error) {
	// Get domain segments for the given domain.
	segments := strings.Split(domain, ".")
	length := len(segments)
	if length < 3 {
		err := errors.New("invalid domain")
		log.Error().Err(err).Msg("Failed to fetch zone ID")
		return nil, err
	}

	// We assume that the zone name is equal to the top-level
	// domain name, except for the dot being replaced by a hyphen.
	zoneID := strings.Join(segments[length-3:length-1], "-")

	// Load credentials file content.
	credentialsBytes, err := ioutil.ReadFile(authFile)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read credentials file")
		return nil, err
	}

	// Parse google credentials.
	var credentials *GoogleCredentials
	if err := json.Unmarshal(credentialsBytes, &credentials); err != nil {
		log.Error().Err(err).Msg("Failed to parse credentials file")
		return nil, err
	}

	// Fetch the project ID.
	projectID := credentials.ProjectID
	if projectID == "" {
		err := errors.New("missing configuration key: project_id")
		log.Error().Err(err).Msg("Failed to fetch project ID")
		return nil, err
	}

	// Create a new service to manage the DNS records. Note that we load the
	// entire file credential file again. This is done as passing the credentials
	// to the service will cause issues with the requested OAuth scopes.
	service, err := dns.NewService(context.TODO(), option.WithCredentialsFile(authFile))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Google Cloud DNS service")
		return nil, err
	}

	return &GoogleProvider{
		ProjectID: projectID,
		ZoneID:    zoneID,
		service:   service,
	}, nil
}

// Reconcile configures the DNS record for the given domain.
func (p *GoogleProvider) Reconcile(domain string, recordType RecordType, ip net.IP) error {
	// Get the IP family.
	ipFamily := 0
	if recordType == RecordTypeA {
		ipFamily = 4
	}
	if recordType == RecordTypeAAAA {
		ipFamily = 6
	}

	// Skip update if the desired IP is empty.
	if ip == nil {
		log.Warn().Msgf("Skipping update: No IPv%d connectivity for record type %s", ipFamily, recordType)
		return nil
	}

	// Define the desired record set.
	targetState := &dns.ResourceRecordSet{
		Name:    domain,
		Type:    string(recordType),
		Rrdatas: []string{ip.String()},
		Ttl:     300,
	}

	// Get the record set for the given domain.
	record, err := p.service.ResourceRecordSets.Get(p.ProjectID, p.ZoneID, domain, string(recordType)).Do()
	if err != nil {
		if e, ok := err.(*googleapi.Error); ok && e.Code == http.StatusNotFound {
			// Create the record, because it does not exist yet.
			_, err := p.service.ResourceRecordSets.Create(p.ProjectID, p.ZoneID, targetState).Do()
			if err != nil {
				log.Error().Err(err).Msg("Failed to create record")
				return err
			}
			log.Info().Msgf("Created record: %s %s %s\n", domain, recordType, ip)
			return nil
		}
		log.Error().Err(err).Msg("Failed to read record")
	}

	// Check if an update may be performed. We will only
	// update the record if it has one a single IP address.
	if len(record.Rrdatas) != 1 {
		log.Warn().Msgf("Skipping update: Multiple IP addresses for record type %s", recordType)
		return nil
	}

	// Check if an update is necessary.
	if record.Rrdatas[0] == ip.String() {
		log.Info().Msgf("Skipping update: %s %s %s", domain, recordType, ip)
		return nil
	}

	// Update the record.
	_, err = p.service.ResourceRecordSets.Patch(p.ProjectID, p.ZoneID, domain, string(recordType), targetState).Do()
	if err != nil {
		log.Error().Err(err).Msg("Failed to update record")
		return err
	}
	log.Info().Msgf("Updated record: %s %s %s\n", domain, recordType, ip)
	return nil
}
