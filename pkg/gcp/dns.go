package gcp

import (
	dns "github.com/pulumi/pulumi-google-native/sdk/go/google/dns/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	ExportZones = "zones"

	ZoneIntric       = "intric-dk"
	ZoneMykilio      = "mykil-io"
	ZoneNicklasfrahm = "nicklasfrahm-xyz"
)

type Zone struct {
	Name        string
	Domain      string
	Description string
}

func StackDNS() pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {
		provider, err := Provider(ctx)
		if err != nil {
			return err
		}

		zones := []Zone{
			{
				Name:        ZoneNicklasfrahm,
				Description: "Nicklas Frahm's personal domain",
				Domain:      "nicklasfrahm.xyz",
			},
			{
				Name:        ZoneIntric,
				Description: "Intric Denmark startup",
				Domain:      "intric.dk",
			},
			{
				Name:        ZoneMykilio,
				Description: "Mykilio project",
				Domain:      "mykil.io",
			},
		}

		for _, zone := range zones {
			managedZone, err := dns.NewManagedZone(ctx, zone.Name, &dns.ManagedZoneArgs{
				Name:    pulumi.StringPtr(zone.Name),
				DnsName: pulumi.StringPtr(zone.Domain),
			}, pulumi.Provider(provider))
			if err != nil {
				return err
			}

			ctx.Export(zone.Name, managedZone.ID().ToIDPtrOutput())
		}

		return nil
	}
}
