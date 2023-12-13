package dns

import (
	"fmt"

	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	// ComponentTypeGithubPages is the ID of the component type.
	ComponentTypeGithubPages = "nicklasfrahm:dns:GithubPages"
)

var (
	// githubPagesIPs contains the IPv4 address of the Github pages servers.
	githubPagesIPs = []string{
		"185.199.108.153",
		"185.199.109.153",
		"185.199.110.153",
		"185.199.111.153",
	}
)

// GithubPages creates DNS records for a Github pages site.
type GithubPages struct {
	pulumi.ResourceState
}

// NewGithubPages configures DNS for a Github pages site.
func NewGithubPages(ctx *pulumi.Context, name string, zone *cloudflare.Zone, args *RecordSpec, opts ...pulumi.ResourceOption) (*GithubPages, error) {
	component := &GithubPages{}
	if err := ctx.RegisterComponentResource(ComponentTypeGithubPages, name, component, opts...); err != nil {
		return nil, err
	}

	if args.GithubPages.Org == "" {
		return nil, fmt.Errorf("%s: failed to find required argument: org", ComponentTypeGithubPages)
	}
	githubPagesSite := fmt.Sprintf("%s.github.io", args.GithubPages.Org)
	isApex := args.Name == "@"

	metadata, err := newMetadataString()
	if err != nil {
		return nil, err
	}

	// Create www record.
	wwwName := fmt.Sprintf("www.%s", args.Name)
	if isApex {
		wwwName = "www"
	}
	_, err = cloudflare.NewRecord(ctx, fmt.Sprintf("%s-r.record-www", name), &cloudflare.RecordArgs{
		ZoneId:  zone.ID(),
		Name:    pulumi.String(wwwName),
		Type:    pulumi.String("CNAME"),
		Value:   pulumi.String(githubPagesSite),
		Comment: pulumi.String(metadata),
	}, pulumi.Parent(component))
	if err != nil {
		return nil, err
	}

	// Create A records for an apex domain or a CNAME record for a subdomain.
	if isApex {
		for _, ip := range githubPagesIPs {
			_, err := cloudflare.NewRecord(ctx, fmt.Sprintf("%s-r.record-%s", name, ip), &cloudflare.RecordArgs{
				ZoneId:  zone.ID(),
				Name:    pulumi.String(args.Name),
				Type:    pulumi.String("A"),
				Value:   pulumi.String(ip),
				Comment: pulumi.String(metadata),
			}, pulumi.Parent(component))
			if err != nil {
				return nil, err
			}
		}
	} else {
		_, err := cloudflare.NewRecord(ctx, fmt.Sprintf("%s-r.record-cname", name), &cloudflare.RecordArgs{
			ZoneId:  zone.ID(),
			Name:    pulumi.String(args.Name),
			Type:    pulumi.String("CNAME"),
			Value:   pulumi.String(githubPagesSite),
			Comment: pulumi.String(metadata),
		}, pulumi.Parent(component))
		if err != nil {
			return nil, err
		}
	}

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{}); err != nil {
		return nil, err
	}

	return component, nil
}
