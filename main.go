package main

import (
	"errors"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nicklasfrahm/infrastructure/pkg/dns"
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
		zones := []dns.Zone{
			dns.NewZone(ZoneNicklasfrahm, "nicklasfrahm.xyz", "Nicklas Frahm's personal domain"),
			dns.NewZone(ZoneIntric, "intric.dk", "Intric Denmark startup"),
			dns.NewZone(ZoneMykilio, "mykil.io", "Mykilio project"),
		}

		// Define GitHub organizations and repositories.
		orgs := []*github.OrganizationConfig{
			github.NewOrganizationConfig("nicklasfrahm", []github.RepositoryConfig{
				github.NewRepositoryConfig("adp", "The source code for the exercises during the ADP course.", nil),
				github.NewRepositoryConfig("archivist", "Process large amounts of files and folders.", nil),
				github.NewRepositoryConfig("distributed-charging", "A proof of concept for a distributed charging infrastructure.", nil),
				github.NewRepositoryConfig("esd", "Code snippets for the ESD course.", nil),
				github.NewRepositoryConfig("file-secret-action", "A GitHub Action to upload a file as a GitHub Actions Secret.", &github.RepositoryExtensions{
					Topics: []string{"github", "secret", "file", "action"},
				}),
				github.NewRepositoryConfig("indesy-robot", "The robot for an indoor delivery system.", &github.RepositoryExtensions{
					Topics: []string{"robot", "system", "indoor", "delivery", "indesy"},
				}),
				github.NewRepositoryConfig("indesy-server", "The server for an indoor delivery system.", &github.RepositoryExtensions{
					Topics: []string{"system", "server", "indoor", "delivery", "indesy"},
				}),
				github.NewRepositoryConfig("indesy-webclient", "The web application for an indoor delivery system.", &github.RepositoryExtensions{
					Topics: []string{"system", "webclient", "indoor", "delivery", "indesy"},
				}),
				// github.NewRepositoryConfig("infrastructure", "421646445"),
				// github.NewRepositoryConfig("labman", "332178426"),
				// github.NewRepositoryConfig("mykilio", "363624141"),
				// github.NewRepositoryConfig("nicklasfrahm", "366965590"),
				// github.NewRepositoryConfig("node-iot-server", "103060200"),
				// github.NewRepositoryConfig("odance", "359234052"),
				// github.NewRepositoryConfig("rtos", "147308829"),
				// github.NewRepositoryConfig("rts", "255901708"),
				// github.NewRepositoryConfig("scp-action", "367666163"),
				// github.NewRepositoryConfig("ses", "222951163"),
				// github.NewRepositoryConfig("showcases", "417350673"),
				// github.NewRepositoryConfig("sonderborg-smart-zero-hack-19-mock-box", "220659897"),
			}),
			github.NewOrganizationConfig("mykilio", []github.RepositoryConfig{
				github.NewRepositoryConfig("docs", "A collection of documentation for all components of the Mykilio ecosystem.", &github.RepositoryExtensions{
					HomepageUrl: "https://docs.mykil.io",
					Topics:      []string{"documentation", "microservices", "mykilio"},
					PagesSource: "docs",
				}),
				github.NewRepositoryConfig("mykilio.go", "The primary SDK for integrating with the ecosystem.", &github.RepositoryExtensions{
					Topics: []string{"go", "microservices", "sdk", "mykilio"},
				}),
			}),
			github.NewOrganizationConfig("intric", []github.RepositoryConfig{
				github.NewRepositoryConfig("services", "", nil),
			}),
		}

		// Define all available stacks.
		stacks := map[string]pulumi.RunFunc{
			StackDNS:    dns.StackDNS(zones),
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
