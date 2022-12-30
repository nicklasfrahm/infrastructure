package stacks

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var Scopes = map[string]pulumi.RunFunc{
	sharedScope: sharedRunFunc,
	dnsScope:    dnsRunFunc,
}
