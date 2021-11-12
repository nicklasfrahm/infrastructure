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
	Repositories []Repository
}

type Organization struct {
	Name     string
	Context  *pulumi.Context
	Provider *github.Provider
}

func NewOrganizationConfig(name string, repos []Repository) *OrganizationConfig {
	return &OrganizationConfig{Name: name, Repositories: repos}
}

func NewOrganization(ctx *pulumi.Context, name string) (*Organization, error) {
	// DEBUG: PAT value
	pat := os.Getenv("PERSONAL_ACCESS_TOKEN")
	fmt.Println(len(pat))

	id := fmt.Sprintf("github-%s", name)
	provider, err := github.NewProvider(ctx, id, &github.ProviderArgs{
		Owner: pulumi.StringPtr(name),
		Token: pulumi.StringPtr(pat),
	})
	if err != nil {
		return nil, err
	}

	return &Organization{
		Name:     name,
		Context:  ctx,
		Provider: provider,
	}, nil
}

type RepositoryConfig struct {
	Name string
}

type Repository struct {
	Name         string
	Organization *Organization
	Repository   *github.Repository
}

func (org *Organization) NewRepository(name string) (*Repository, error) {

	id := fmt.Sprintf("%s-%s", org.Name, name)
	repo, err := github.NewRepository(org.Context, id, &github.RepositoryArgs{
		Name:                pulumi.String(name),
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
	}, pulumi.Provider(org.Provider), pulumi.Parent(org.Provider))
	if err != nil {
		return nil, err
	}

	return &Repository{
		Name:         name,
		Organization: org,
		Repository:   repo,
	}, nil
}
