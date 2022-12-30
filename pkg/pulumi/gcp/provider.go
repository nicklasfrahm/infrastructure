package gcp

import (
	"fmt"

	"github.com/pulumi/pulumi-google-native/sdk/go/google"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// GetProvider returns a pre-configured Google Cloud provider.
func GetProvider(ctx *pulumi.Context) (*google.Provider, error) {
	projectKey := "google-native:project"
	project, hasProject := ctx.GetConfig(projectKey)
	if !hasProject {
		return nil, fmt.Errorf("missing configuration key: %s", projectKey)
	}

	regionKey := "google-native:region"
	region, hasRegion := ctx.GetConfig(regionKey)
	if !hasRegion {
		return nil, fmt.Errorf("missing configuration key: %s", regionKey)
	}

	provider, err := google.NewProvider(ctx, "provider-freeTier", &google.ProviderArgs{
		Project: pulumi.StringPtr(project),
		Region:  pulumi.StringPtr(region),
	})
	if err != nil {
		return nil, err
	}

	return provider, nil
}
