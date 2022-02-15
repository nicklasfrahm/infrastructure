package main

import (
	"os"

	"github.com/spf13/cobra"
)

var help bool
var authFile string

var rootCmd = &cobra.Command{
	Use:   "dnsadm",
	Short: "CLI to automate DNS management",
	Long: `A command line interface to automate DNS management tasks. Currently
it only uses Google Cloud DNS, but this may change in the future.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if help {
			cmd.Help()
			os.Exit(0)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(0)
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "display help for command")
	rootCmd.PersistentFlags().StringVarP(&authFile, "auth-file", "a", "/etc/secrets/credentials.json", "set cloud credentials path")
}
