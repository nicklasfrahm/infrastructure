package dns

import (
	"log"
	"strings"

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

func NewZone(domain string, description string) Zone {
	return Zone{
		Name:        strings.ReplaceAll(domain, ".", "-"),
		Domain:      domain,
		Description: description,
	}
}

func Stack() pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {
		zones := []Zone{
			NewZone("intric.dk", "Intric Denmark"),
			NewZone("nicklasfrahm.xyz", "Nicklas Frahm"),
			NewZone("mykil.io", "Mykilio Project"),
		}

		for _, zone := range zones {
			managedZone, err := dns.NewManagedZone(ctx, zone.Name, &dns.ManagedZoneArgs{
				Name:    pulumi.StringPtr(zone.Name),
				DnsName: pulumi.StringPtr(zone.Domain),
			})
			if err != nil {
				return err
			}

			log.Println(managedZone.ID().ToStringOutput())
		}

		// TODO: Export managed zones somehow.
		// ctx.Export(ExportZones, "test")
		return nil
	}
}
