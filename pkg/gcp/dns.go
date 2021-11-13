package gcp

import (
	"fmt"

	dns "github.com/pulumi/pulumi-google-native/sdk/go/google/dns/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// For more information about how DNSSEC works, consider the link below.
// https://www.cloudflare.com/dns/dnssec/how-dnssec-works/

const (
	// Enable DNSSEC to prevent DNS spoofing.
	DnsSecState = dns.ManagedZoneDnsSecConfigStateOn
	// Use NSEC to reduce cryptographic complexity under
	// the assumption that our applications are safe and
	// that zone-walking does not pose a threat.
	DnsSecNonExistence = dns.ManagedZoneDnsSecConfigNonExistenceNsec
	// Use the maximum number of bits, where the increment is 64,
	// the minimum is 512 and the maximum is 1024.
	DnsSecKeyLength = 1024
)

type Zone struct {
	Name        string
	Domain      string
	Description string
}

func NewZone(name string, domain string, description string) Zone {
	return Zone{
		Name:        name,
		Domain:      domain,
		Description: description,
	}
}

func StackDNS(zones []Zone) pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {
		provider, err := Provider(ctx)
		if err != nil {
			return err
		}

		for _, zone := range zones {
			managedZone, err := dns.NewManagedZone(ctx, zone.Name, &dns.ManagedZoneArgs{
				Name:        pulumi.String(zone.Name),
				DnsName:     pulumi.String(fmt.Sprintf("%s.", zone.Domain)),
				Description: pulumi.String(zone.Description),
				// Official REST API reference for configuration options:
				// https://cloud.google.com/dns/docs/reference/v1/managedZones
				// Information about advanced DNSSEC setups:
				// https://cloud.google.com/dns/docs/dnssec-advanced
				DnssecConfig: dns.ManagedZoneDnsSecConfigArgs{
					State:        DnsSecState,
					NonExistence: DnsSecNonExistence,
					DefaultKeySpecs: dns.DnsKeySpecArray{
						dns.DnsKeySpecArgs{
							Algorithm: dns.DnsKeySpecAlgorithmEcdsap384sha384,
							KeyLength: pulumi.IntPtr(DnsSecKeyLength),
							KeyType:   dns.DnsKeySpecKeyTypeZoneSigning,
						},
						dns.DnsKeySpecArgs{
							Algorithm: dns.DnsKeySpecAlgorithmEcdsap384sha384,
							KeyLength: pulumi.IntPtr(DnsSecKeyLength),
							KeyType:   dns.DnsKeySpecKeyTypeKeySigning,
						},
					},
				},
				Visibility: dns.ManagedZoneVisibilityPublic,
			}, pulumi.Provider(provider), pulumi.Parent(provider))
			if err != nil {
				return err
			}

			// Export the ID of the zone.
			ctx.Export(fmt.Sprintf("managedzone.id/%s", zone.Domain), managedZone.ID())
		}

		return nil
	}
}
