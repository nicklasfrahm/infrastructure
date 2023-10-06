package main

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"
var help bool

var rootCmd = &cobra.Command{
	Use:   "hwctl",
	Short: "A command line tool to manage physical servers.",
	Long: `   _                 _   _
  | |____      _____| |_| |
  | '_ \ \ /\ / / __| __| |
  | | | \ V  V / (__| |_| |
  |_| |_|\_/\_/ \___|\__|_|

A command line tool to manage a physical server fleet.
The tool provides helpers to manage an inventory and
credentials for the servers.`,
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
	Version:      version,
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "display help for command")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
