package metal

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "metal",
	Short: `Manage appliances and their lifecycle`,
	Long: `This subcommand allows you to provision
and manage appliances with a pre-installed
operating system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Command.AddCommand(bootstrapCmd)
}
