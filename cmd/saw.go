package cmd

import (
	"github.com/spf13/cobra"
)

// SawCommand is the main top-level command
var SawCommand = &cobra.Command{
	Use:   "saw <command>",
	Short: "A fast, multipurpose tool for AWS CloudWatch Logs",
	Long:  "Saw is a fast, multipurpose tool for AWS CloudWatch Logs.",
	Example: `  saw groups
  saw streams production
  saw watch production`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

func init() {
	SawCommand.AddCommand(groupsCommand)
	SawCommand.AddCommand(streamsCommand)
	SawCommand.AddCommand(versionCommand)
	SawCommand.AddCommand(watchCommand)
	SawCommand.AddCommand(getCommand)
}
