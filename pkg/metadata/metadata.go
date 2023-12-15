package metadata

import "encoding/json"

// Manager describes which service is administering the infrastructure.
type Manager string

const (
	// ManagerPulumi indicates that pulumi is managing the infrastructure.
	ManagerPulumi Manager = "pulumi"
)

// Metadata is a data structure that contains metadata about a piece of infrastructure.
type Metadata struct {
	// Manager is the name of the service administering the infrastructure.
	Manager Manager `json:"manager"`
}

// JSON returns the JSON representation of the metadata.
func (m *Metadata) JSON() ([]byte, error) {
	return json.Marshal(m)
}
