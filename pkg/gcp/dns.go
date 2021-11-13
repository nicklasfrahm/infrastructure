package gcp

import (
	"fmt"

	dns "github.com/pulumi/pulumi-google-native/sdk/go/google/dns/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	ExportZones = "zones"
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
