package zone

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"

	"github.com/nicklasfrahm/infrastructure/pkg/kraut/zone"
)

var (
	zoneConfig = zone.Zone{
		Router: &zone.ZoneRouter{},
	}
	configFile string
)

var upCmd = &cobra.Command{
	Use:   "up <host>",
	Short: "Bootstrap a new availability zone",
	Long: `This command will bootstrap a new zone by connecting
to the specified IP and setting up a k3s cluster on
the host that will then set up the required services
for managing the lifecycle of the zone.

To manage a zone, the CLI needs credentials for the
DNS provider that is used to manage the DNS records
for the zone. These credentials can only be provided
via the environment variable DNS_PROVIDER_CREDENTIAL
and DNS_PROVIDER or via a ".env" file in the current
working directory.`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"host"},
	ValidArgs:  []string{"host"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// This should be safe because of the ExactArgs(1) constraint,
		// but we still need to check it to avoid panics.
		if len(args) != 1 {
			return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
		}
		host := args[0]

		fmt.Printf("zone: %+v\n", zoneConfig)
		fmt.Printf("rout: %+v\n", zoneConfig.Router)
		fmt.Printf("prov: %+v\n", os.Getenv("DNS_PROVIDER"))

		if err := zone.Up(host, &zoneConfig); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	upCmd.Flags().StringVarP(&zoneConfig.Name, "name", "n", "", "name of the zone")
	upCmd.Flags().StringVarP(&zoneConfig.Domain, "domain", "d", "", "domain that will contain the DNS records for the zone")
	upCmd.Flags().StringVarP(&zoneConfig.Router.Hostname, "hostname", "H", "", "hostname of the router serving the zone")
	upCmd.Flags().IPVarP(&zoneConfig.Router.ID, "router-id", "r", nil, "IPv4 address of the router serving the zone")
	upCmd.Flags().Uint32VarP(&zoneConfig.Router.ASN, "asn", "a", 0, "autonomous system number of the zone")
	upCmd.Flags().StringVarP(&configFile, "config", "c", "", "path to the configuration file")
}
