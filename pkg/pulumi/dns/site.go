package dns

import (
	"fmt"

	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	// SiteComponentType is the ID of the component type.
	SiteComponentType = "nicklasfrahm:dns:Site"
)

// GithubPages creates DNS records for a Github pages site.
type Site struct {
	pulumi.ResourceState
}

// NewSite configures DNS for a Github pages site.
func NewSite(ctx *pulumi.Context, name string, zone *cloudflare.Zone, args *RecordSpec, opts ...pulumi.ResourceOption) (*Site, error) {
	component := &Site{}
	if err := ctx.RegisterComponentResource(SiteComponentType, name, component, opts...); err != nil {
		return nil, err
	}

	if args.Site.Router == "" {
		return nil, fmt.Errorf("%s: failed to find required argument: router", SiteComponentType)
	}

	_, err := cloudflare.NewRecord(ctx, fmt.Sprintf("%s-r.record-base", name), &cloudflare.RecordArgs{
		ZoneId: zone.ID(),
		Name:   pulumi.String(args.Name),
		Type:   pulumi.String("CNAME"),
		Value:  pulumi.String(args.Site.Router),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	_, err = cloudflare.NewRecord(ctx, fmt.Sprintf("%s-r.record-wildcard", name), &cloudflare.RecordArgs{
		ZoneId: zone.ID(),
		Name:   pulumi.Sprintf("*.%s", args.Name),
		Type:   pulumi.String("CNAME"),
		Value:  pulumi.String(args.Site.Router),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{}); err != nil {
		return nil, err
	}

	return component, nil
}
