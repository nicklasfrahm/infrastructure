package dns

import (
	"github.com/nicklasfrahm/infrastructure/pkg/metadata"
)

// newMetadataString generates a string containing
// metadata about the infrastructure.
func newMetadataString() (string, error) {
	meta := &metadata.Metadata{
		Manager: metadata.ManagerPulumi,
	}

	metaJSON, err := meta.JSON()
	if err != nil {
		return "", err
	}

	return string(metaJSON), nil
}
