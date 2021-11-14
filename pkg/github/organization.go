package github

import (
	"fmt"
	"os"
	"strings"

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
	Name        string
	Description string
	Private     bool
	HomepageUrl string
	Topics      []string
	// The format for the pages source is `<branch>:/path/to/folder`.
	// This implicity enables GitHub pages.
	PagesSource string
}

type RepositoryExtensions struct {
	Private     bool
	HomepageUrl string
	Topics      []string
	// The format for the pages source is `<branch>:/path/to/folder`.
	// This implicity enables GitHub pages.
	PagesSource string
}

type Repository struct {
	Name         string
	Organization *Organization
	Repository   *github.Repository
}

func NewRepositoryConfig(
	name string,
	description string,
	extensions *RepositoryExtensions,
) RepositoryConfig {
	extras := new(RepositoryExtensions)
	if extensions != nil {
		if extensions.Topics == nil {
			extensions.Topics = []string{}
		}

		*extras = *extensions
	}

	return RepositoryConfig{
		Name:        name,
		Description: description,
		Private:     extras.Private,
		HomepageUrl: extras.HomepageUrl,
		Topics:      extras.Topics,
		PagesSource: extras.PagesSource,
	}
}

func (org *Organization) NewRepository(config *RepositoryConfig) (*Repository, error) {
	// Configuration for the repository.
	repository := github.RepositoryArgs{
		Name:                pulumi.String(config.Name),
		Description:         pulumi.String(config.Description),
		Visibility:          pulumi.String("public"),
		HomepageUrl:         pulumi.String(config.HomepageUrl),
		AllowAutoMerge:      pulumi.Bool(true),
		AllowMergeCommit:    pulumi.Bool(true),
		AllowRebaseMerge:    pulumi.Bool(true),
		AllowSquashMerge:    pulumi.Bool(true),
		DeleteBranchOnMerge: pulumi.Bool(true),
		HasWiki:             pulumi.Bool(false),
		HasIssues:           pulumi.Bool(true),
		HasProjects:         pulumi.Bool(true),
		VulnerabilityAlerts: pulumi.Bool(true),
		HasDownloads:        pulumi.Bool(true),
		AutoInit:            pulumi.Bool(false),
		Topics:              pulumi.ToStringArray(config.Topics),
		// ArchiveOnDestroy:    pulumi.Bool(false),
	}

	// Overwrite visibility.
	if config.Private {
		repository.Visibility = pulumi.String("private")
	}

	// Configure GitHub pages.
	if config.PagesSource != "" {
		// The format for the pages source is `<branch>:/path/to/folder`.
		sourceParameters := strings.Split(config.PagesSource, ":")
		source := github.RepositoryPagesSourceArgs{
			Branch: pulumi.String(sourceParameters[0]),
		}

		// Overwrite the repository root as the default path if provided.
		if len(sourceParameters) == 2 {
			source.Path = pulumi.String(sourceParameters[1])
		}

		// Get domain name.
		domainSegments := strings.Split(config.HomepageUrl, "://")

		repository.Pages = github.RepositoryPagesArgs{
			Cname:  pulumi.String(domainSegments[len(domainSegments)-1]),
			Source: source,
		}
	}

	id := fmt.Sprintf("%s-%s", org.Name, config.Name)
	repo, err := github.NewRepository(org.Context, id, &repository, pulumi.Provider(org.Provider), pulumi.Parent(org.Provider), pulumi.Import(pulumi.ID(config.Name)))
	if err != nil {
		return nil, err
	}

	return &Repository{
		Name:         config.Name,
		Organization: org,
		Repository:   repo,
	}, nil
}
