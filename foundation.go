package main

import (
	"fmt"
	"os"

	"github.com/nicklasfrahm/infrastructure/pkg/pulumi/dns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"
)

const (
	// StackFoundation is the name of the stack that deploys DNS,
	// Kubernetes at the network edge and underlay networking.
	StackFoundation = "foundation"
	// dnsSpecPath is the path to the DNS specification.
	dnsSpecPath = "deploy/dns.yaml"
)

// NewFoundationStack deploys DNS, Kubernetes at the network edge
// and underlay networking.
func NewFoundationStack(ctx *pulumi.Context) error {
	if err := configureDNS(ctx); err != nil {
		return err
	}

	return nil
}

// configureDNS loads the DNS specification and configures DNS.
func configureDNS(ctx *pulumi.Context) error {
	dnsSpecBytes, err := os.ReadFile(dnsSpecPath)
	if err != nil {
		return err
	}

	var dnsSpec dns.Spec
	if err := yaml.Unmarshal(dnsSpecBytes, &dnsSpec); err != nil {
		return err
	}

	for i := 0; i < len(dnsSpec.Zones); i++ {
		zone := &dnsSpec.Zones[i]

		_, err := dns.NewZone(ctx, fmt.Sprintf("%s-c.zone-%s", StackFoundation, zone.Name), zone)
		if err != nil {
			return err
		}
	}

	return nil
}
