package dns

import (
	"fmt"

	"github.com/pulumi/pulumi-cloudflare/sdk/v5/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	// AComponentType is the ID of the component type.
	AComponentType = "nicklasfrahm:dns:A"
)

// A creates A DNS records for the given hostname.
type A struct {
	pulumi.ResourceState
}

// NewA configures A DNS records for the given hostname.
func NewA(ctx *pulumi.Context, name string, zone *cloudflare.Zone, args *RecordSpec, opts ...pulumi.ResourceOption) (*A, error) {
	component := &A{}
	if err := ctx.RegisterComponentResource(AComponentType, name, component, opts...); err != nil {
		return nil, err
	}

	if len(args.Values) == 0 {
		return nil, fmt.Errorf("%s: failed to find required argument: values", AComponentType)
	}

	metadata, err := newMetadataString()
	if err != nil {
		return nil, err
	}

	for _, value := range args.Values {

		_, err := cloudflare.NewRecord(ctx, fmt.Sprintf("%s-r.record-%s", name, value), &cloudflare.RecordArgs{
			ZoneId:  zone.ID(),
			Name:    pulumi.String(args.Name),
			Type:    pulumi.String("A"),
			Value:   pulumi.String(value),
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
