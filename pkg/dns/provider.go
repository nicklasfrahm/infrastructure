package dns

import (
	"fmt"

	"github.com/pulumi/pulumi-google-native/sdk/go/google"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	DefaultRegion = "europe-north1"
	Project       = "nicklasfrahm"
)

var gcp *google.Provider

func ProviderGCP(ctx *pulumi.Context) (*google.Provider, error) {
	if gcp != nil {
		return gcp, nil
	}

	provider, err := google.NewProvider(ctx, fmt.Sprintf("gcp-%s", Project), &google.ProviderArgs{
		Project: pulumi.String(Project),
		Region:  pulumi.String(DefaultRegion),
	})
	if err != nil {
		return nil, err
	}

	gcp = provider

	return gcp, nil
}
