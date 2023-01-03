package stacks

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nicklasfrahm/infrastructure/pkg/pulumi/gcp"
)

var dnsScope = "dns"

func dnsRunFunc(ctx *pulumi.Context) error {
	name := dnsScope

	provider, err := gcp.NewProvider(ctx, name)
	if err != nil {
		return err
	}

	dynDNSFunction, err := gcp.NewPublicFunction(ctx, fmt.Sprintf("%s-c.function-dynDNS", name), &gcp.PublicFunctionArgs{
		Name:       "dyndns",
		EntryPoint: "UpdateDNSRecord",
		Runtime:    "go119",
	}, pulumi.Provider(provider))
	if err != nil {
		return err
	}
	_ = dynDNSFunction

	return nil
}
