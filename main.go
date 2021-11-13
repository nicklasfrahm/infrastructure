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
			github.NewOrganizationConfig("nicklasfrahm", []github.RepositoryConfig{
				github.NewRepositoryConfig("adp", "120659570"),
				github.NewRepositoryConfig("archivist", "352418616"),
				github.NewRepositoryConfig("distributed-charging", "329911865"),
				github.NewRepositoryConfig("esd", "219478595"),
				github.NewRepositoryConfig("file-secret-action", "367789870"),
				github.NewRepositoryConfig("indesy-robot", "106269089"),
				github.NewRepositoryConfig("indesy-server", "106265257"),
				github.NewRepositoryConfig("indesy-webclient", "104065906"),
				github.NewRepositoryConfig("infrastructure", "421646445"),
				github.NewRepositoryConfig("labman", "332178426"),
				github.NewRepositoryConfig("mykilio", "363624141"),
				github.NewRepositoryConfig("nicklasfrahm", "366965590"),
				github.NewRepositoryConfig("node-iot-server", "103060200"),
				github.NewRepositoryConfig("odance", "359234052"),
				github.NewRepositoryConfig("rtos", "147308829"),
				github.NewRepositoryConfig("rts", "255901708"),
				github.NewRepositoryConfig("scp-action", "367666163"),
				github.NewRepositoryConfig("ses", "222951163"),
				github.NewRepositoryConfig("showcases", "417350673"),
				github.NewRepositoryConfig("sonderborg-smart-zero-hack-19-mock-box", "220659897"),
			}),
			github.NewOrganizationConfig("mykilio", []github.RepositoryConfig{
				github.NewRepositoryConfig("docs", "423123012"),
				github.NewRepositoryConfig("mykilio.go", "426119679"),
			}),
			github.NewOrganizationConfig("intric", []github.RepositoryConfig{
				// github.NewRepositoryConfig("services", "427475385"),
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
