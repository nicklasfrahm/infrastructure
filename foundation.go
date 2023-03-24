package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// FoundationStack is the name of the stack that deploys DNS,
// Kubernetes at the network edge and underlay networking.
const StackFoundation = "foundation"

// NewFoundationStack deploys DNS, Kubernetes at the network edge
// and underlay networking.
func NewFoundationStack(ctx *pulumi.Context) error {
	return nil
}
