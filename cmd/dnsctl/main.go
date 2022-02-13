package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

const (
	// EndpointIPv4 is an endpoint to fetch the public IPv4 address.
	EndpointIPv4 = "https://ip4.seeip.org"
	// EndpointIPv6 is an endpoint to fetch the public IPv6 address.
	EndpointIPv6 = "https://ip6.seeip.org"
	// CommandName is the name of the command.
	CommandName = "dnsctl"
)

var usage = `A dynamic DNS client to update or create Google Cloud DNS
records. By default it will load the credentials from the
/etc/secrets/gcp.json JSON file. This behaviour can be
overwritten by setting the CREDENTIALS_FILE environment
variable. The service account in question must have the
DNS Admin role.

Usage:
  dnsctl <domain>`

// RecordType is the type of DNS record.
type RecordType string

const (
	// RecordTypeA is the A record type.
	RecordTypeA RecordType = "A"
	// RecordTypeAAAA is the AAAA record type.
	RecordTypeAAAA RecordType = "AAAA"
)

func main() {
	// Configure logger.
	log.SetFlags(log.LUTC | log.Ltime | log.Ldate)

	// Parse command line arguments.
	fqdn := ParseCommandLine()

	// Fetch Google Cloud credentials and use them
	// to create a service for the Cloud DNS API.
	provider := NewGoogleProvider()

	// Get the desired state.
	ipv4, err := PublicIP(EndpointIPv4)
	if err != nil {
		log.Fatalf("❌\tFailed to get IPv4 address: %v", err)
	}
	ipv6, err := PublicIP(EndpointIPv6)
	if err != nil {
		log.Fatalf("❌\tFailed to get IPv6 address: %v", err)
	}

	// Reconcile the records.
	provider.Reconcile(fqdn, RecordTypeA, ipv4)
	provider.Reconcile(fqdn, RecordTypeAAAA, ipv6)
}

// ParseConmandLine parses the command line arguments
// and returns the FQDN.
func ParseCommandLine() string {
	// Check if the command line contains the correct number of arguments.
	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	// Strip the program name from the command line.
	args := os.Args[1:]

	// Ensure that the domain name has the
	// canonical format with the trailing dot.
	if !strings.HasSuffix(args[0], ".") {
		args[0] = args[0] + "."
	}

	return args[0]
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
	Reconcile(domain, recordType, desired string)
}

// GoogleCredentials describes the structure
// of the Google Cloud credentials.
type GoogleCredentials struct {
	ProjectID string `json:"project_id"`
}

// GoogleProvider implements the Provider interface
// for Google Cloud DNS.
type GoogleProvider struct {
	credentials *GoogleCredentials
	service     *dns.Service
}

// NewGoogleProvider creates a new GoogleProvider.
func NewGoogleProvider() *GoogleProvider {
	// Fetch credential file location from environment variable.
	credentialsFile := os.Getenv("CREDENTIALS_FILE")
	if credentialsFile == "" {
		credentialsFile = "/etc/secrets/gcp.json"
	}

	// Load credentials file content.
	credentialsBytes, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("❌\tFailed to read credentials file: %v", err)
	}

	// Parse google credentials.
	var credentials *GoogleCredentials
	if err := json.Unmarshal(credentialsBytes, &credentials); err != nil {
		log.Fatalf("❌\tFailed to parse credentials secret file: %v", err)
	}

	// Create a new service to manage the DNS records. Note that we load the
	// entire file credential file again. This is done as passing the credentials
	// to the service will cause issues with the requested OAuth scopes.
	service, err := dns.NewService(context.TODO(), option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Fatalf("❌\tFailed to create DNS service: %v", err)
	}

	return &GoogleProvider{
		credentials: credentials,
		service:     service,
	}
}

// projectID returns the project ID for the given credentials.
func (p *GoogleProvider) projectID() string {
	// Fetch the project ID.
	projectID := p.credentials.ProjectID
	if projectID == "" {
		log.Fatalln("⚠️\tMust provide valid project ID in credentials")
	}

	return projectID
}

// zoneID returns the zone ID for the given domain.
func (p *GoogleProvider) zoneID(domain string) string {
	// Get domain segments for the given domain.
	segments := strings.Split(domain, ".")
	length := len(segments)
	if length < 3 {
		log.Fatalln("⚠️\tMust provide valid top-level domain")
	}

	// We assume that the zone name is equal to the top-level
	// domain name, except for the dot being replaced by a hyphen.
	return strings.Join(segments[length-3:length-1], "-")
}

// Reconcile configures the DNS record for the given domain.
func (p *GoogleProvider) Reconcile(domain string, recordType RecordType, ip net.IP) {
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
		log.Printf("🟡\tSkipping update: No IPv%d connectivity\n", ipFamily)
		return
	}

	// Define the desired record set.
	targetState := &dns.ResourceRecordSet{
		Name:    domain,
		Type:    string(recordType),
		Rrdatas: []string{ip.String()},
		Ttl:     300,
	}

	// Get the record set for the given domain.
	record, err := p.service.ResourceRecordSets.Get(p.projectID(), p.zoneID(domain), domain, string(recordType)).Do()
	if err != nil {
		if e, ok := err.(*googleapi.Error); ok && e.Code == http.StatusNotFound {
			// Create the record, because it does not exist yet.
			_, err := p.service.ResourceRecordSets.Create(p.projectID(), p.zoneID(domain), targetState).Do()
			if err != nil {
				log.Fatalf("❌\tFailed to create record: %v", err)
			}
			log.Printf("🟢\tCreated record: %s %s %s\n", domain, recordType, ip)
			return
		}
		log.Fatalf("❌\tFailed to get record: %v", err)
	}

	// Check if an update may be performed. We will only
	// update the record if it has one a single IP address.
	if len(record.Rrdatas) != 1 {
		log.Printf("🟡\tSkipping update: %s record must have single IPv%d address\n", recordType, ipFamily)
	}

	// Check if an update is necessary.
	if record.Rrdatas[0] == ip.String() {
		log.Printf("🟡\tSkipping update: %s record already up-to-date\n", recordType)
		return
	}

	// Update the record.
	_, err = p.service.ResourceRecordSets.Patch(p.projectID(), p.zoneID(domain), domain, string(recordType), targetState).Do()
	if err != nil {
		log.Fatalf("❌\tFailed to update record: %v", err)
	}
	log.Printf("🟢\tUpdated record: %s %s %s\n", domain, recordType, ip)
}
