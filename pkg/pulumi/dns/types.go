package dns

const (
	// RecordTypeGithubPages is the type of a GitHub Pages DNS record.
	RecordTypeGithubPages = "GITHUB_PAGES"
	// RecordTypeSite is the type of a DNS record for a site.
	RecordTypeSite = "SITE"
	// RecordTypeCNAME is the type of a CNAME DNS record.
	RecordTypeCNAME = "CNAME"
	// RecordTypeA is the type of an A DNS record.
	RecordTypeA = "A"
)

// GitHubPagesRecordSpec is a data structure that describes a GitHub Pages DNS record.
type GitHubPagesRecordSpec struct {
	// Org is the name of the GitHub organization or user that owns the repository.
	Org string `yaml:"org"`
}

// SiteRecordSpec is a data structure that describes a DNS record for a site.
type SiteRecordSpec struct {
	// Router is the FQDN of the edge router.
	Router string `yaml:"router"`
}

// RecordSpec is a data structure that describes a DNS record.
type RecordSpec struct {
	// Name is the name of the DNS record.
	Name string `yaml:"name" validate:"required"`
	// Type is the type of the DNS record.
	Type string `yaml:"type" validate:"required,oneof=GITHUB_PAGES SITE CNAME"`
	// Values is a list of values for the DNS record.
	Values []string `yaml:"values"`
	// GithubPages is configures the GitHub pages site.
	GithubPages GitHubPagesRecordSpec `yaml:"githubPages"`
	// Site configures a site DNS record.
	Site SiteRecordSpec `yaml:"site"`
}

// ZoneSpec is a data structure that describes a DNS zone.
type ZoneSpec struct {
	// Provider is the name of the DNS provider.
	Provider string `yaml:"provider" validate:"required,oneof=cloudflare"`
	// Name is the name of the DNS zone.
	Name string `yaml:"name" validate:"required"`
	// ID is unique identifier of the DNS zone.
	ID string `yaml:"id"`
	// Records is a list of DNS records.
	Records []RecordSpec `yaml:"records" validate:"required,dive,required"`
}

// Spec is a data structure that describes the DNS configuration.
type Spec struct {
	// Zones is a list of DNS zones.
	Zones []ZoneSpec `yaml:"zones" validate:"required,dive,required"`
}
