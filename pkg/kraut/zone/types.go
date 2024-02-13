package zone

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

var (
	// ErrInvalidHostname is returned when a hostname is invalid.
	ErrInvalidHostname = fmt.Errorf("hostname must be lowercase, alphanumeric, and cannot start or end with a hyphen")
	// ErrRouterConfigRequired is returned when a router configuration is required.
	ErrRouterConfigRequired = fmt.Errorf("router configuration is required")
	// ErrRouterIDUnspecified is returned when a router ID is unspecified.
	ErrRouterIDUnspecified = fmt.Errorf("router ID is unspecified")
	// ErrASNRequired is returned when an ASN is required.
	ErrASNRequired = fmt.Errorf("ASN is required")
	// ErrApexDomainRequired is returned when an apex domain is required.
	ErrApexDomainRequired = fmt.Errorf("apex domain is required")
	// ErrNameRequired is returned when a name is required.
	ErrNameRequired = fmt.Errorf("name is required")
)

// ZoneRouter represents the configuration of a router serving an availability zone.
type ZoneRouter struct {
	// Hostname is the hostname of the router.
	Hostname string `json:"hostname"`
	// ID is the IPv4 address of the router.
	ID net.IP `json:"routerID"`
	// ASN is the autonomous system number of the router.
	ASN uint32 `json:"asn"`
}

// Validate returns an error if the router configuration is not valid.
func (r *ZoneRouter) Validate() error {
	// Ensure that the hostname contains only alphanumeric characters and hyphens,
	// and that it does not start or end with a hyphen.
	hostnameRegex := regexp.MustCompile((`^[a-z0-9]+([a-z0-9-]*[a-z0-9]+)*$`))
	if !hostnameRegex.MatchString(r.Hostname) {
		return ErrInvalidHostname
	}

	if r.ID == nil || r.ID.IsUnspecified() {
		return ErrRouterIDUnspecified
	}

	if r.ASN == 0 {
		return ErrASNRequired
	}

	return nil
}

// Zone represents an availability zone.
type Zone struct {
	// Name is the name of the zone.
	Name string `json:"name"`
	// Domain is the domain that will contain the DNS records for the zone.
	Domain string `json:"domain"`
	// Router is the router serving the zone.
	Router *ZoneRouter `json:"router"`
}

// Validate returns an error if the zone is not valid.
func (z *Zone) Validate() error {
	if z.Name == "" {
		return ErrNameRequired
	}

	if strings.Count(strings.Trim(z.Domain, "."), ".") != 1 {
		return ErrApexDomainRequired
	}

	if z.Router == nil {
		return ErrRouterConfigRequired
	}

	if err := z.Router.Validate(); err != nil {
		return err
	}

	return nil
}
