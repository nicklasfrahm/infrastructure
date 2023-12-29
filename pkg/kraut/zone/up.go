package zone

import (
	"fmt"

	"github.com/nicklasfrahm/infrastructure/pkg/netutil"
)

// Up creates or updates an availability zone.
func Up(host string, zone *Zone) error {
	if err := zone.Validate(); err != nil {
		return err
	}

	fingerprint, err := netutil.GetSSHHostPublicKeyFingerprint(host)
	if err != nil {
		return err
	}

	fmt.Printf("fingerprint detected: %s: %s\n", host, fingerprint)

	return nil
}
