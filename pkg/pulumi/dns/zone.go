package dns

import (
	"fmt"
	"os"

	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	// ZoneComponentType is the ID of the component type.
	ZoneComponentType = "nicklasfrahm:dns:Zone"
	// ZoneProviderCloudflare is the name of the Cloudflare provider.
	ZoneProviderCloudflare = "cloudflare"

	// cloudflareEnvVarAPIToken is the name of the environment variable for the API token.
	cloudflareEnvVarAPIToken = "CLOUDFLARE_API_TOKEN"
	// cloudflareEnvVarAccountID is the name of the environment variable for the API token.
	cloudflareEnvVarAccountID = "CLOUDFLARE_ACCOUNT_ID"
	// cloudflarePlan is the plan for the Cloudflare zones.
	cloudflarePlan = "free"
)

// Zone is a high-level component for a DNS zone.
type Zone struct {
	pulumi.ResourceState
}

// NewZone is a high-level abstraction that creates a DNS zone.
// It also takes care of the provider configuration.
func NewZone(ctx *pulumi.Context, name string, args *ZoneSpec, opts ...pulumi.ResourceOption) (*Zone, error) {
	component := &Zone{}
	if err := ctx.RegisterComponentResource(ZoneComponentType, name, component, opts...); err != nil {
		return nil, err
	}

	if args.Provider == ZoneProviderCloudflare {
		provider, err := cloudflare.NewProvider(ctx, fmt.Sprintf("%s-p.cloudflare", name), &cloudflare.ProviderArgs{
			ApiToken: pulumi.StringPtr(os.Getenv(cloudflareEnvVarAPIToken)),
		}, pulumi.Parent(component))
		if err != nil {
			return nil, err
		}

		zoneOptions := []pulumi.ResourceOption{pulumi.Parent(provider), pulumi.Provider(provider)}
		if args.ID != "" {
			zoneOptions = append(zoneOptions, pulumi.Import(pulumi.ID(args.ID)))
		}

		zone, err := cloudflare.NewZone(ctx, fmt.Sprintf("%s-r.zone", name), &cloudflare.ZoneArgs{
			AccountId: pulumi.String(os.Getenv(cloudflareEnvVarAccountID)),
			Zone:      pulumi.String(args.Name),
			Plan:      pulumi.StringPtr(cloudflarePlan),
		}, zoneOptions...)
		if err != nil {
			return nil, err
		}

		for j := 0; j < len(args.Records); j++ {
			record := &args.Records[j]

			switch record.Type {
			case RecordTypeGithubPages:
				_, err = NewGithubPages(ctx, fmt.Sprintf("%s-c.githubpages-%s", name, record.Name), zone, record, pulumi.Parent(component))
			}
			if err != nil {
				return nil, err
			}
		}
	}

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{}); err != nil {
		return nil, err
	}

	return component, nil
}
