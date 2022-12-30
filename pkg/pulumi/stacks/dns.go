package stacks

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nicklasfrahm/infrastructure/pkg/pulumi/gcp"
)

var dnsScope = "dns"

func dnsRunFunc(ctx *pulumi.Context) error {
	name := dnsScope

	provider, err := gcp.GetProvider(ctx)
	if err != nil {
		return err
	}

	dynDNSFunction, err := gcp.NewPublicFunction(ctx, fmt.Sprintf("%s-c.function-dynDNS", name), &gcp.PublicFunctionArgs{}, pulumi.Provider(provider))
	if err != nil {
		return err
	}
	_ = dynDNSFunction

	return nil
}
