package stacks

import (
	"fmt"

	storage "github.com/pulumi/pulumi-google-native/sdk/go/google/storage/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nicklasfrahm/infrastructure/pkg/pulumi/gcp"
)

var sharedScope = "shared"

func sharedRunFunc(ctx *pulumi.Context) error {
	name := sharedScope

	provider, err := gcp.NewProvider(ctx, name)
	if err != nil {
		return err
	}

	// Create a storage bucket that will contain all the source code for my cloud functions.
	_, err = storage.NewBucket(ctx, fmt.Sprintf("%s-r.bucket-functionSource", name), &storage.BucketArgs{
		Name:         pulumi.Sprintf("%s-functionsource", provider.Project.Elem()).ToStringPtrOutput(),
		Location:     provider.Region,
		LocationType: pulumi.StringPtr("region"),
		IamConfiguration: &storage.BucketIamConfigurationArgs{
			PublicAccessPrevention: pulumi.StringPtr("enforced"),
			UniformBucketLevelAccess: &storage.BucketIamConfigurationUniformBucketLevelAccessArgs{
				Enabled: pulumi.Bool(false),
			},
		},
		Versioning: &storage.BucketVersioningArgs{
			// This is necessary to ensure that Google Cloud Functions can
			// detect changes to the source code to rebuild the function.
			Enabled: pulumi.Bool(true),
		},
	}, pulumi.Provider(provider), pulumi.DeleteBeforeReplace(true))
	if err != nil {
		return err
	}

	return nil
}
