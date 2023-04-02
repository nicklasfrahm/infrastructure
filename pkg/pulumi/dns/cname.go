package dns

import (
	"fmt"

	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	// CNAMEComponentType is the ID of the component type.
	CNAMEComponentType = "nicklasfrahm:dns:CNAME"
)

// CNAME creates CNAME DNS records for the given hostname.
type CNAME struct {
	pulumi.ResourceState
}

// NewCNAME configures CNAME DNS records for the given hostname.
func NewCNAME(ctx *pulumi.Context, name string, zone *cloudflare.Zone, args *RecordSpec, opts ...pulumi.ResourceOption) (*CNAME, error) {
	component := &CNAME{}
	if err := ctx.RegisterComponentResource(CNAMEComponentType, name, component, opts...); err != nil {
		return nil, err
	}

	if len(args.Values) == 0 {
		return nil, fmt.Errorf("%s: failed to find required argument: values", CNAMEComponentType)
	}

	for _, value := range args.Values {
		_, err := cloudflare.NewRecord(ctx, fmt.Sprintf("%s-r.record-%s", name, value), &cloudflare.RecordArgs{
			ZoneId: zone.ID(),
			Name:   pulumi.String(args.Name),
			Type:   pulumi.String("CNAME"),
			Value:  pulumi.String(value),
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
