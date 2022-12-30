package gcp

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// PublicFunctionArgs are the arguments for creating a public function.
type PublicFunctionArgs struct {
}

// PublicFunction is a public function that can be invoked via HTTP.
type PublicFunction struct {
	pulumi.ResourceState
}

// NewPublicFunction creates a public function that can be invoked via HTTP.
func NewPublicFunction(ctx *pulumi.Context, name string, args *PublicFunctionArgs, opts ...pulumi.ResourceOption) (*PublicFunction, error) {
	component := &PublicFunction{}
	err := ctx.RegisterComponentResource("nicklasfrahm:gcp:PublicFunction", name, component, opts...)
	if err != nil {
		return nil, err
	}

	// TODO: Implement bucket object creation and function creation.
	// Reference: https://github.com/pulumi/pulumi-google-native/blob/master/examples/functions-ts/index.ts

	return component, nil

}
