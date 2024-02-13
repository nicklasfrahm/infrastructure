package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/nicklasfrahm/infrastructure/cmd/ic/zone"
)

var version = "dev"
var help bool

var rootCmd = &cobra.Command{
	Use:   "ic",
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
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		os.Exit(1)
		return nil
	},
	Version:      version,
	SilenceUsage: true,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "Print this help")
	rootCmd.AddCommand(zone.Command)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
