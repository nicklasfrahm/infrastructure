package zone

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	zone = Zone{
		Router: &ZoneRouter{},
	}
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
		// TODO: Validate args:
		//       - IP (e.g. 212.112.144.171) [ephemeral, required]
		//       - hostname (e.g. alfa.nicklasfrahm.dev) [config, required]
		//       - zone name (e.g. aar1) [config, required]
		//       - router ID (e.g. 172.31.255.0) [config, required]
		//       - ASN (e.g. 65000) [config, required]
		//       - DNS provider credential [env only, required]
		//       - user account password for local recovery [env only, optional]
		// TODO: Run preflight checks:
		//       - Open ports: TCP:22,80,443,6443,7443
		//       - Open ports: UDP:5800-5810
		// TODO: Perform minimal system configuration:
		//       - Set hostname
		//       - Reset user password (if provided)
		// TODO: Ensure minimal interface configuration:
		//       - IPv4 on loopback
		//       - Identify WAN interface and name it WAN
		//       - DHCP on all interfaces (if not configured)
		// TODO: Install or upgrade k3s
		// TODO: Install or upgrade kraut

		host := args[0]
		fingerprint, err := GetSSHHostPublicKeyFingerprint(host)
		if err != nil {
			return err
		}

		fmt.Printf("fingerprint detected: %s: %s\n", host, fingerprint)

		return nil
	},
}

func init() {
	upCmd.Flags().StringVarP(&zone.Name, "name", "n", "", "name of the zone")
	upCmd.MarkFlagRequired("name")
	upCmd.Flags().StringVarP(&zone.Domain, "domain", "d", "", "domain that will contain the DNS records for the zone")
	upCmd.MarkFlagRequired("domain")
	upCmd.Flags().StringVarP(&zone.Router.Hostname, "hostname", "H", "", "hostname of the router serving the zone")
	upCmd.MarkFlagRequired("hostname")
	upCmd.Flags().IPVarP(&zone.Router.ID, "router-id", "r", nil, "IPv4 address of the router serving the zone")
	upCmd.MarkFlagRequired("router-id")
	upCmd.Flags().Uint32VarP(&zone.Router.ASN, "asn", "a", 0, "autonomous system number of the zone")
	upCmd.MarkFlagRequired("asn")
}

func GetSSHHostPublicKeyFingerprint(host string) (string, error) {
	// Read the private key file.
	// TODO: Avoid hardcoding this.
	key, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/id_ed25519")
	if err != nil {
		return "", fmt.Errorf("failed to read private key file: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// This will be written by the host key callback.
	var fingerprint string

	config := &ssh.ClientConfig{
		User: os.Getenv("USER"),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// Create a SHA256 fingerprint of the host public key.
			fingerprint = ssh.FingerprintSHA256(key)

			return nil
		}),
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(host, "22"), config)
	if err != nil {
		return "", fmt.Errorf("failed to establish SSH connection: %v", err)
	}
	defer client.Close()

	return fingerprint, nil
}
