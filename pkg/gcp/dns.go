package gcp

import (
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
				DnsName:     pulumi.String(zone.Domain),
				Description: pulumi.String(zone.Description),
			}, pulumi.Provider(provider))
			if err != nil {
				return err
			}

			ctx.Export(zone.Name, managedZone.ID().ToIDPtrOutput())
		}

		return nil
	}
}
