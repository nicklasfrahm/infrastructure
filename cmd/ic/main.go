package main

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var version = "dev"
var help bool

var rootCmd = &cobra.Command{
	Use:   "ic <host>",
	Short: "A CLI to manage infrastructure",
	Long: `   _
  (_) ___
  | |/ __|
  | | (__
  |_|\___|

ic is a CLI to manage infrastructure. It provides
a variety of commands to manage different stages
of the infrastructure lifecycle.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if help {
			cmd.Help()
			os.Exit(0)
		}
	},
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"host"},
	ValidArgs:  []string{"host"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			cmd.Help()
			os.Exit(1)
		}

		host := args[0]
		fingerprint, err := GetSSHHostPublicKeyFingerprint(host)
		if err != nil {
			return err
		}

		fmt.Printf("fingerprint detected: %s: %s\n", host, fingerprint)

		return nil
	},
	Version:      version,
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "Print this help")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
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
