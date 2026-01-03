package cmd

import "github.com/spf13/cobra"

var RootCmd = &cobra.Command{
	Use:   "ntmux",
	Short: "Yet another tmux wrapper",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
