package dns

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ZoneComponentType is the ID of the component type.
const ZoneComponentType = "nicklasfrahm:dns:Zone"

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

	return component, nil
}
