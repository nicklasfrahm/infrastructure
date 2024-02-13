package sshx

import (
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/nicklasfrahm/infrastructure/pkg/netutil"
)

// NewClient creates a new SSH client with all available SSH keys
// in the user's home directory via the default port 22.
func NewClient(host string) (*ssh.Client, error) {
	// Get the user's home directory
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %s", err)
	}

	// Create SSH auth methods using the files in the user's home directory.
	authMethods, err := createAuthMethods(filepath.Join(usr.HomeDir, ".ssh"))
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH auth methods: %s", err)
	}

	// Create an SSH client config with the specified authentication methods
	config := &ssh.ClientConfig{
		User:            usr.Username,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(host, fmt.Sprint(netutil.PortSSH)), config)
	if err != nil {
		return nil, fmt.Errorf("failed to establish SSH connection: %s", err)
	}

	return client, nil
}

// createAuthMethods creates SSH authentication methods using private key files
// in the specified directory.
func createAuthMethods(dir string) ([]ssh.AuthMethod, error) {
	var authMethods []ssh.AuthMethod

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range files {
		// Avoid directories and public key files.
		if entry.IsDir() || strings.HasSuffix(entry.Name(), ".pub") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		key, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			// Skip files that are not private keys.
			continue
		}

		fmt.Println(path)

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("failed to find private key files: %s", dir)
	}

	return authMethods, nil
}

// GetSSHHostPublicKeyFingerprint returns the SHA256 fingerprint of the host.
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
