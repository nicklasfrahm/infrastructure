package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var stacks = map[string]pulumi.RunFunc{
	StackFoundation: NewFoundationStack,
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stack := ctx.Stack()

		// Check if the stack has a run function that defines
		// which infrastructure is to be deployed.
		runFunc := stacks[stack]
		if runFunc == nil {
			return fmt.Errorf("failed to find stack: %s", stack)
		}

		return runFunc(ctx)
	})
}
