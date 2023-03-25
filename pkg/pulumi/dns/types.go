package dns

// RecordSpec is a data structure that describes a DNS record.
type RecordSpec struct {
	// Name is the name of the DNS record.
	Name string `yaml:"name"`
	// Type is the type of the DNS record.
	Type string `yaml:"type"`
	// Values is a list of values for the DNS record.
	Values []string `yaml:"values"`
}

// ZoneSpec is a data structure that describes a DNS zone.
type ZoneSpec struct {
	// Provider is the name of the DNS provider.
	Provider string `yaml:"provider" validate:"required,oneof=cloudflare"`
	// Name is the name of the DNS zone.
	Name string `yaml:"name" validate:"required"`
	// ID is unique identifier of the DNS zone.
	ID string `yaml:"id"`
}

// Spec is a data structure that describes the DNS configuration.
type Spec struct {
	// Zones is a list of DNS zones.
	Zones []ZoneSpec `yaml:"zones" validate:"required,dive,required"`
}
