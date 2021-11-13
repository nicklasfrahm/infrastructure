package github

import (
	"fmt"
	"os"

	"github.com/pulumi/pulumi-github/sdk/v4/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	VisibilityPublic = "public"
)

type OrganizationConfig struct {
	Name         string
	Repositories []RepositoryConfig
}

type Organization struct {
	Name     string
	Context  *pulumi.Context
	Provider *github.Provider
}

func NewOrganizationConfig(name string, repos []RepositoryConfig) *OrganizationConfig {
	return &OrganizationConfig{Name: name, Repositories: repos}
}

func NewOrganization(ctx *pulumi.Context, config *OrganizationConfig) (*Organization, error) {
	id := fmt.Sprintf("github-%s", config.Name)
	provider, err := github.NewProvider(ctx, id, &github.ProviderArgs{
		Owner: pulumi.StringPtr(config.Name),
		Token: pulumi.StringPtr(os.Getenv("PERSONAL_ACCESS_TOKEN")),
	})
	if err != nil {
		return nil, err
	}

	return &Organization{
		Name:     config.Name,
		Context:  ctx,
		Provider: provider,
	}, nil
}

type RepositoryConfig struct {
	ID   string
	Name string
}

type Repository struct {
	Name         string
	Organization *Organization
	Repository   *github.Repository
}

func NewRepositoryConfig(name string, id string) RepositoryConfig {
	return RepositoryConfig{
		ID:   id,
		Name: name,
	}
}

func (org *Organization) NewRepository(config *RepositoryConfig) (*Repository, error) {
	id := fmt.Sprintf("%s-%s", org.Name, config.Name)
	ref := fmt.Sprintf("%s/%s", org.Name, config.Name)
	repo, err := github.NewRepository(org.Context, id, &github.RepositoryArgs{
		Name:                pulumi.String(config.Name),
		Visibility:          pulumi.String(VisibilityPublic),
		AllowAutoMerge:      pulumi.Bool(true),
		AllowMergeCommit:    pulumi.Bool(true),
		AllowRebaseMerge:    pulumi.Bool(true),
		AllowSquashMerge:    pulumi.Bool(true),
		ArchiveOnDestroy:    pulumi.Bool(true),
		AutoInit:            pulumi.Bool(false),
		DeleteBranchOnMerge: pulumi.Bool(true),
		HasIssues:           pulumi.Bool(true),
		HasProjects:         pulumi.Bool(true),
		HasWiki:             pulumi.Bool(false),
		VulnerabilityAlerts: pulumi.Bool(true),
	}, pulumi.Provider(org.Provider), pulumi.Parent(org.Provider), pulumi.Import(pulumi.ID(ref)))
	if err != nil {
		return nil, err
	}

	return &Repository{
		Name:         config.Name,
		Organization: org,
		Repository:   repo,
	}, nil
}
