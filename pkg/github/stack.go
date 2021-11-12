package github

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func StackGitHub(configs []*OrganizationConfig) pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {
		for _, config := range configs {
			// Create a new organization.
			organization, err := NewOrganization(ctx, config.Name)
			if err != nil {
				return err
			}

			for _, repo := range config.Repositories {
				_, err := organization.NewRepository(repo.Name)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}
}
