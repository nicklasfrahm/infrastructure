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

type RepositoryOptions struct {
	Private     bool
	HomepageUrl string
	Topics      []string
	// The format for the pages source is `<branch>:/path/to/folder`.
	// This implicity enables GitHub pages.
	PagesSource string
	Template    string
}

type RepositoryConfig struct {
	Name        string
	Description string
	Options     *RepositoryOptions
}

type Repository struct {
	Name         string
	Organization *Organization
	Repository   *github.Repository
}

func NewRepositoryConfig(
	name string,
	description string,
	options *RepositoryOptions,
) RepositoryConfig {
	if options == nil {
		options = new(RepositoryOptions)
	}

	if options.Topics == nil {
		options.Topics = make([]string, 0)
	}

	return RepositoryConfig{
		Name:        name,
		Description: description,
		Options:     options,
	}
}

func (org *Organization) NewRepository(config *RepositoryConfig) (*Repository, error) {
	// Configuration for the repository.
	repository := github.RepositoryArgs{
		Name:                pulumi.String(config.Name),
		Description:         pulumi.String(config.Description),
		Visibility:          pulumi.String("public"),
		HomepageUrl:         pulumi.String(config.Options.HomepageUrl),
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
		Topics:              pulumi.ToStringArray(config.Options.Topics),
		// ArchiveOnDestroy:    pulumi.Bool(false),
	}

	// Overwrite visibility.
	if config.Options.Private {
		repository.Visibility = pulumi.String("private")
	}

	// Configure GitHub pages.
	if config.Options.PagesSource != "" {
		// The format for the pages source is `<branch>:/path/to/folder`.
		sourceParameters := strings.Split(config.Options.PagesSource, ":")
		source := github.RepositoryPagesSourceArgs{
			Branch: pulumi.String(sourceParameters[0]),
		}

		// Overwrite the repository root as the default path if provided.
		if len(sourceParameters) == 2 {
			source.Path = pulumi.String(sourceParameters[1])
		}

		// Get domain name.
		domainSegments := strings.Split(config.Options.HomepageUrl, "://")

		repository.Pages = github.RepositoryPagesArgs{
			Cname:  pulumi.String(domainSegments[len(domainSegments)-1]),
			Source: source,
		}
	}

	// Configure template reference.
	if config.Options.Template != "" {
		chunks := strings.Split(config.Options.Template, "/")
		owner := "nicklasfrahm"
		repo := chunks[0]
		if len(chunks) >= 2 {
			owner = chunks[0]
			repo = chunks[1]
		}
		repository.Template = github.RepositoryTemplateArgs{
			Owner:      pulumi.String(owner),
			Repository: pulumi.String(repo),
		}
	}

	id := fmt.Sprintf("%s-%s", org.Name, config.Name)
	options := []pulumi.ResourceOption{
		pulumi.Provider(org.Provider),
		pulumi.Parent(org.Provider),
		pulumi.Import(pulumi.ID(config.Name)),
	}
	repo, err := github.NewRepository(org.Context, id, &repository, options...)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Name:         config.Name,
		Organization: org,
		Repository:   repo,
	}, nil
}
