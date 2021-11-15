package kubernetes

import (
	"strings"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	continents = map[string]string{
		"eur": "Europe",
	}
	countries = map[string]string{
		"den": "Denmark",
		"ger": "Germany",
	}
	cities = map[string]string{
		"sdb": "Sonderborg",
	}
)

type Region struct {
	Name      string
	Continent string
	Country   string
	City      string
	Domain    string
}

func NewRegion(name string) *Region {
	chunks := strings.Split(name, "-")

	// Make sure we have all chunks we need.
	if len(chunks) < 3 {
		for len(chunks) < 3 {
			chunks = append(chunks, "")
		}
	}

	return &Region{
		Continent: continents[chunks[0]],
		Country:   countries[chunks[1]],
		City:      cities[chunks[2]],
	}
}

func Stack() pulumi.RunFunc {
	return func(ctx *pulumi.Context) error {
		return nil
	}
}

// Define infrastructure regions. The naming scheme consists of:
// - three character continent: eur
// - three character country: den, ger
// - three character city: sdb, flb
// clusters := []*Cluster{
//   dns.NewRegion("eur-den-sdb"),
// }
