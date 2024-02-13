package netutil

import (
	"fmt"
	"net"
	"time"
)

// ProbeStatus represents the status of a port.
type ProbeStatus string

const (
	// ProbeStatusOpen indicates that a port is open.
	ProbeStatusOpen ProbeStatus = "open"
	// ProbeStatusClosed indicates that a port is closed.
	ProbeStatusClosed ProbeStatus = "closed"
	// ProbeStatusFiltered indicates that a port is filtered.
	ProbeStatusFiltered ProbeStatus = "filtered"
)

const (
	// PortSSH is the default SSH port.
	PortSSH = 22
)

// ProbeTCP checks if a port is open on a host by
// connecting to it with a timeout of 1 second.
func ProbeTCP(host string, port int) ProbeStatus {
	status := ProbeStatusOpen

	if _, err := net.DialTimeout("tcp", net.JoinHostPort(host, fmt.Sprint(port)), time.Second); err != nil {
		status = ProbeStatusFiltered

		// Check if the error is a I/O timeout error.
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			status = ProbeStatusClosed
		}
	}

	return status
}
