package main

import (
	"errors"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nicklasfrahm/infrastructure/pkg/dns"
)

const (
	StackDNS = "dns"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Define all available stacks.
		stacks := map[string]pulumi.RunFunc{
			StackDNS: dns.Stack(),
		}

		// Load the corresponding stack and verify its existence.
		run := stacks[ctx.Stack()]
		if run == nil {
			return errors.New("pulumi: unknown stack")
		}

		// Run stack to reconcile the resources.
		return run(ctx)
	})
}
