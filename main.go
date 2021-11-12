package main

import (
	"errors"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nicklasfrahm/infrastructure/pkg/gcp"
	"github.com/nicklasfrahm/infrastructure/pkg/github"
)

const (
	StackDNS    = "dns"
	StackGitHub = "github"

	ZoneIntric       = "intric-dk"
	ZoneMykilio      = "mykil-io"
	ZoneNicklasfrahm = "nicklasfrahm-xyz"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Define DNS zones.
		zones := []gcp.Zone{
			gcp.NewZone(ZoneNicklasfrahm, "nicklasfrahm.xyz", "Nicklas Frahm's personal domain"),
			gcp.NewZone(ZoneIntric, "intric.dk", "Intric Denmark startup"),
			gcp.NewZone(ZoneMykilio, "mykil.io", "Mykilio project"),
		}

		// Define GitHub organizations and repositories.
		orgs := []*github.OrganizationConfig{
			github.NewOrganizationConfig("nicklasfrahm", []github.Repository{
				{Name: "infrastructure"},
			}),
			github.NewOrganizationConfig("mykilio", []github.Repository{
				{Name: "docs"},
				{Name: "mykilio.go"},
			}),
			github.NewOrganizationConfig("intric", []github.Repository{
				{Name: "services"},
			}),
		}

		// Define all available stacks.
		stacks := map[string]pulumi.RunFunc{
			StackDNS:    gcp.StackDNS(zones),
			StackGitHub: github.StackGitHub(orgs),
		}

		// Load the corresponding stack and verify its existence.
		run := stacks[ctx.Stack()]
		if run == nil {
			return errors.New("pulumi: unknown stack")
		}

		// Run stack to reconcile the resources.
		return run(ctx)
	})
}
