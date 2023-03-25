package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// stacks maps stack names to run functions that define
// which infrastructure is to be deployed.
var stacks = map[string]pulumi.RunFunc{
	StackFoundation: NewFoundationStack,
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stack := ctx.Stack()

		// Check if the stack has an implementation.
		runFunc := stacks[stack]
		if runFunc == nil {
			return fmt.Errorf("failed to find stack: %s", stack)
		}

		return runFunc(ctx)
	})
}
