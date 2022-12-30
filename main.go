package main

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nicklasfrahm/infrastructure/pkg/pulumi/stacks"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackID := ctx.Stack()

		stackInfo := strings.Split(stackID, "-")

		// The stack consists of two parts, the logical
		// scope name and the environment name.
		if len(stackInfo) != 2 {
			return fmt.Errorf("invalid stack name does not follow format: <scope>-<env>")
		}

		// Check if the scope has a registered run fuction
		// that defines which infrastructure to deploy.
		runFunc := stacks.Scopes[stackInfo[0]]
		if runFunc == nil {
			return fmt.Errorf("invalid scope: %s", stackInfo[0])
		}

		return runFunc(ctx)
	})
}
