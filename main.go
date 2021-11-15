package main

import (
	"errors"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nicklasfrahm/infrastructure/pkg/dns"
	"github.com/nicklasfrahm/infrastructure/pkg/github"
	"github.com/nicklasfrahm/infrastructure/pkg/kubernetes"
)

const (
	StackDNS        = "dns"
	StackGitHub     = "github"
	StackKubernetes = "kubernetes"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Define DNS zones.
		zones := []*dns.Zone{
			dns.NewZone("nicklasfrahm.xyz", "Nicklas Frahm's personal domain"),
			dns.NewZone("mykil.io", "Mykilio project"),
			dns.NewZone("intric.dk", "Intric Denmark startup"),
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
				github.NewRepositoryConfig("infrastructure", "A repository to automate the configuration of my local infrastructure.", &github.RepositoryExtensions{
					Topics: []string{"infrastructure-as-code", "hybrid-cloud"},
				}),
				github.NewRepositoryConfig("labman", "A cloud-based management solution for laboratories and workshops.", nil),
				github.NewRepositoryConfig("mykilio", "Mykilio is the proposal for a new Living Standard with the goal to reimagine infrastructure management and monitoring. It aims to be scalable, lightweight, extensible and secure by applying IoT principles to the datacenter and homelabs.", &github.RepositoryExtensions{
					Topics:      []string{"infrastructure", "monitoring", "baremetal", "mykilio"},
					HomepageUrl: "https://mykil.io",
					PagesSource: "docs",
				}),
				github.NewRepositoryConfig("nicklasfrahm", "An overview of my current projects and ideas.", nil),
				github.NewRepositoryConfig("node-iot-server", "A small server to make a serial device addressable over UDP.", &github.RepositoryExtensions{
					Topics: []string{"iot", "serial", "server", "udp"},
				}),
				github.NewRepositoryConfig("odance", "O!Dance is a dancing school project in Sonderborg, Denmark.", &github.RepositoryExtensions{
					Topics:      []string{"website", "school", "dancing"},
					HomepageUrl: "https://odance.dk",
					PagesSource: "gh-pages",
				}),
				github.NewRepositoryConfig("rtos", "The source code for the exercises during the RTOS course.", nil),
				github.NewRepositoryConfig("rts", "Real-time systems", nil),
				github.NewRepositoryConfig("scp-action", "A Github Action to upload and download files via SCP.", &github.RepositoryExtensions{
					Topics: []string{"github", "ssh", "scp", "action"},
				}),
				github.NewRepositoryConfig("ses", "Software for Embedded Systems Course Source Code", nil),
				github.NewRepositoryConfig("showcases", "A repository containing showcases for my programming and architecture skills.", &github.RepositoryExtensions{
					Topics:      []string{"go", "microservices", "nats", "event-driven"},
					HomepageUrl: "https://nicklasfrahm.xyz",
					PagesSource: "gh-pages",
				}),
				github.NewRepositoryConfig("sonderborg-smart-zero-hack-19-mock-box", "", nil),
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
			StackDNS:        dns.Stack(zones),
			StackKubernetes: kubernetes.Stack(),
			StackGitHub:     github.Stack(orgs),
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
