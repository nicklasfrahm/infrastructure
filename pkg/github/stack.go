package github

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Stack(configs []*OrganizationConfig) pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {
		for _, org := range configs {
			// Create a new organization.
			organization, err := NewOrganization(ctx, org)
			if err != nil {
				return err
			}

			// Create repositories for the organization.
			for _, repo := range org.Repositories {
				_, err := organization.NewRepository(&repo)
				if err != nil {
					return err
				}

				// TODO: Create DNS records for GitHub pages, if zone is managed by Pulumi.
			}
		}

		return nil
	}
}
