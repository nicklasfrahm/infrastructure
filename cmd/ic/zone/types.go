package zone

import (
	"fmt"
	"net"
	"regexp"
)

// TODO: Move this to the `pkg` directory.

var (
	// ErrInvalidHostname is returned when a hostname is invalid.
	ErrInvalidHostname = fmt.Errorf("hostname must be lowercase, alphanumeric, and cannot start or end with a hyphen")
)

// Zone represents an availability zone.
type Zone struct {
	// Name is the name of the zone.
	Name string `json:"name"`
	// Hostname is the hostname of the router serving the zone.
	Hostname string `json:"hostname"`
	// RouterID is the IPv4 address of the router serving the zone.
	RouterID net.IP `json:"routerID"`
	// ASN is the autonomous system number of the zone.
	ASN int `json:"asn"`
	// Domain is the domain that will contain the DNS records for the zone.
	Domain string `json:"comain"`
}

// Validate returns an error if the zone is not valid.
func (z *Zone) Validate() error {
	// Ensure that the hostname contains only alphanumeric characters and hyphens,
	// and that it does not start or end with a hyphen.
	hostnameRegex := regexp.MustCompile((`^[a-z0-9]+([a-z0-9-]*[a-z0-9]+)*$`))
	if !hostnameRegex.MatchString(z.Hostname) {
		return ErrInvalidHostname
	}

	return nil
}
