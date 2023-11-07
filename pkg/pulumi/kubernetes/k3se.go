package dns

import (
	"github.com/go-playground/validator/v10"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	// ComponentTypeK3se is the ID of the component type.
	ComponentTypeK3se = "nicklasfrahm:kubernetes:K3se"
)

// K3se manages a k3s cluster using k3se.
type K3se struct {
	pulumi.ResourceState
}

// K3seArgs are the arguments for creating a k3s cluster using k3se.
type K3seArgs struct {
	// Name is the name of the cluster.
	Name string `validate:"required,alphanum"`
}

// NewK3se configures CNAME DNS records for the given hostname.
func NewK3se(ctx *pulumi.Context, name string, args *K3seArgs, opts ...pulumi.ResourceOption) (*K3se, error) {
	component := &K3se{}
	if err := ctx.RegisterComponentResource(ComponentTypeK3se, name, component, opts...); err != nil {
		return nil, err
	}

	if err := validator.New().Struct(args); err != nil {
		return nil, err
	}

	// TODO: Implement logic here.

	if err := ctx.RegisterResourceOutputs(component, pulumi.Map{}); err != nil {
		return nil, err
	}

	return component, nil
}
