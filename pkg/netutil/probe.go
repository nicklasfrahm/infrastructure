package netutil

import (
	"fmt"
	"net"
	"time"
)

// ProbeTCP checks if a port is open on a host by
// connecting to it with a timeout of 1 second.
func ProbeTCP(host string, port int) bool {
	if _, err := net.DialTimeout("tcp", net.JoinHostPort(host, fmt.Sprint(port)), time.Second); err != nil {
		return false
	}

	return true
}
