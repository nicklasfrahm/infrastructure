package gcp

import (
	"fmt"

	"github.com/pulumi/pulumi-google-native/sdk/go/google"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	projectKey = "google-native:project"
	regionKey  = "google-native:region"
)

// NewProvider returns a pre-configured Google Cloud provider.
func NewProvider(ctx *pulumi.Context, name string, opts ...pulumi.ResourceOption) (*google.Provider, error) {
	project, err := GetProject(ctx)
	if err != nil {
		return nil, err
	}

	region, err := GetRegion(ctx)
	if err != nil {
		return nil, err
	}

	provider, err := google.NewProvider(ctx, fmt.Sprintf("%s-p.googleNative-freeTier", name), &google.ProviderArgs{
		Project: pulumi.StringPtr(project),
		Region:  pulumi.StringPtr(region),
	}, opts...)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// GetProject returns the GCP project ID.
func GetProject(ctx *pulumi.Context) (string, error) {
	project, hasProject := ctx.GetConfig(projectKey)
	if !hasProject {
		return "", fmt.Errorf("missing configuration key: %s", projectKey)
	}

	return project, nil
}

// GetRegion returns the GCP region.
func GetRegion(ctx *pulumi.Context) (string, error) {
	region, hasRegion := ctx.GetConfig(regionKey)
	if !hasRegion {
		return "", fmt.Errorf("missing configuration key: %s", regionKey)
	}

	return region, nil
}
