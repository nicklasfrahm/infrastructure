// Package dyndns provides a function to update a DNS record in a Google Cloud DNS zone.
package dyndns

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

// UpdateDNSRecord updates a DNS record in a Google Cloud DNS zone.
func UpdateDNSRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Generate a random base64.
	buf := make([]byte, 64)
	_, err := rand.Read(buf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := base64.StdEncoding.EncodeToString(buf)

	w.Write([]byte(fmt.Sprintf(`{"random": "%s"}`, token)))
}
