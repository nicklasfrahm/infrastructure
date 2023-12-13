package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"

	"github.com/nicklasfrahm/infrastructure/pkg/pulumi/dns"
	"github.com/nicklasfrahm/infrastructure/pkg/pulumi/kubernetes"
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
	zoneResources, err := configureDNS(ctx)
	if err != nil {
		return err
	}

	if err := configureClusters(ctx, pulumi.DependsOn(zoneResources)); err != nil {
		return err
	}

	return nil
}

// configureDNS loads the DNS specification and configures DNS.
func configureDNS(ctx *pulumi.Context, opts ...pulumi.ResourceOption) ([]pulumi.Resource, error) {
	dnsSpecBytes, err := os.ReadFile(dnsSpecPath)
	if err != nil {
		return nil, err
	}

	var dnsSpec dns.Spec
	if err := yaml.Unmarshal(dnsSpecBytes, &dnsSpec); err != nil {
		return nil, err
	}

	if err := validator.New().Struct(dnsSpec); err != nil {
		return nil, err
	}

	resources := make([]pulumi.Resource, len(dnsSpec.Zones))
	for i := 0; i < len(dnsSpec.Zones); i++ {
		zone := &dnsSpec.Zones[i]

		resource, err := dns.NewZone(ctx, fmt.Sprintf("%s-c.zone-%s", StackFoundation, zone.Name), zone, opts...)
		if err != nil {
			return nil, err
		}

		resources[i] = resource
	}

	return resources, nil
}

// configureEdgeClusters creates k3s clusters at the network edge.
func configureClusters(ctx *pulumi.Context, opts ...pulumi.ResourceOption) error {
	entries, err := os.ReadDir(kubernetes.K3seSpecDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		clusterName := strings.TrimSuffix(entry.Name(), ".yaml")
		cluster, err := kubernetes.NewK3se(ctx, fmt.Sprintf("%s-c.k3se-%s", StackFoundation, clusterName), &kubernetes.K3seArgs{
			Name: clusterName,
		}, opts...)
		if err != nil {
			return err
		}

		if clusterName == "charlie" {
			_, err := corev1.NewNamespace(ctx, fmt.Sprintf("%s-%s-r.namespace-test", StackFoundation, clusterName), &corev1.NamespaceArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Name: pulumi.String("pulumi-test"),
				},
			}, pulumi.Parent(cluster), pulumi.Provider(cluster.Provider))
			if err != nil {
				return err
			}
		}

		// TODO: Remove this.
		cluster.Server.ApplyT(func(server string) string {
			if server != "" {
				pulumi.Printf("%s: %s\n", clusterName, server)
			}

			return server
		})
	}

	return nil
}
