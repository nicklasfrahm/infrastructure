package gcp

import (
	"fmt"

	cloudFunctions "github.com/pulumi/pulumi-google-native/sdk/go/google/cloudfunctions/v1"
	storage "github.com/pulumi/pulumi-google-native/sdk/go/google/storage/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// PublicFunctionArgs are the arguments for creating a public function.
type PublicFunctionArgs struct {
	Name       string
	EntryPoint string
	Runtime    string
}

// PublicFunction is a public function that can be invoked via HTTP.
type PublicFunction struct {
	pulumi.ResourceState
}

// NewPublicFunction creates a public function that can be invoked via HTTP.
// The implementation is inspired by the TypeScript example on Github:
// https://github.com/pulumi/pulumi-google-native/blob/master/examples/functions-ts/index.ts
func NewPublicFunction(ctx *pulumi.Context, name string, args *PublicFunctionArgs, opts ...pulumi.ResourceOption) (*PublicFunction, error) {
	component := &PublicFunction{}
	err := ctx.RegisterComponentResource("nicklasfrahm:gcp:PublicFunction", name, component, opts...)
	if err != nil {
		return nil, err
	}

	project, err := GetProject(ctx)
	if err != nil {
		return nil, err
	}

	region, err := GetRegion(ctx)
	if err != nil {
		return nil, err
	}

	sourceLocation := fmt.Sprintf("./functions/%s", args.Name)

	bucketObject, err := storage.NewBucketObject(ctx, fmt.Sprintf("%s-r.bucketObject-source", name), &storage.BucketObjectArgs{
		Name:   pulumi.StringPtr(args.Name),
		Bucket: pulumi.Sprintf("%s-functionsource", project),
		Source: pulumi.NewAssetArchive(map[string]interface{}{
			".": pulumi.NewFileArchive(sourceLocation),
		}),
	}, pulumi.Parent(component), pulumi.ReplaceOnChanges([]string{"source"}))
	if err != nil {
		return nil, err
	}

	function, err := cloudFunctions.NewFunction(ctx, fmt.Sprintf("%s-r.function-http", name), &cloudFunctions.FunctionArgs{
		Name:              pulumi.Sprintf("projects/%s/locations/%s/functions/%s", project, region, args.Name),
		SourceArchiveUrl:  pulumi.Sprintf("gs://%s/%s", bucketObject.Bucket, bucketObject.Name),
		EntryPoint:        pulumi.StringPtr(args.EntryPoint),
		Runtime:           pulumi.StringPtr(args.Runtime),
		Timeout:           pulumi.StringPtr("60s"),
		AvailableMemoryMb: pulumi.IntPtr(128),
		HttpsTrigger: &cloudFunctions.HttpsTriggerArgs{
			SecurityLevel: cloudFunctions.HttpsTriggerSecurityLevelSecureAlways,
		},
		IngressSettings: cloudFunctions.FunctionIngressSettingsAllowAll,
		// This is a hack to trigger a re-deployment when the source code changes.
		BuildEnvironmentVariables: pulumi.StringMap{
			"SOURCE_HASH": bucketObject.Crc32c,
		},
	}, pulumi.Parent(component), pulumi.DependsOn([]pulumi.Resource{bucketObject}))
	if err != nil {
		return nil, err
	}

	_, err = cloudFunctions.NewFunctionIamPolicy(ctx, fmt.Sprintf("%s-r.functionIAMPolicy-allUsers", name), &cloudFunctions.FunctionIamPolicyArgs{
		FunctionId: pulumi.String(args.Name),
		Bindings: cloudFunctions.BindingArray{
			&cloudFunctions.BindingArgs{
				Role: pulumi.String("roles/cloudfunctions.invoker"),
				Members: pulumi.StringArray{
					pulumi.String("allUsers"),
				},
			},
		},
	}, pulumi.Parent(component), pulumi.DependsOn([]pulumi.Resource{function}))
	if err != nil {
		return nil, err
	}

	return component, nil

}
