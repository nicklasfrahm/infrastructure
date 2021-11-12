package gcp

import (
	gcp "github.com/pulumi/pulumi-google-native/sdk/go/google"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	DefaultRegion = "europe-north1"
)

var provider *gcp.Provider

func Provider(ctx *pulumi.Context) (*gcp.Provider, error) {
	if provider != nil {
		return provider, nil
	}

	p, err := gcp.NewProvider(ctx, "gcp", &gcp.ProviderArgs{
		Region: pulumi.String(DefaultRegion),
	})
	if err != nil {
		return nil, err
	}

	provider = p

	return provider, nil
}
