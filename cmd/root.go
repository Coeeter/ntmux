package cmd

import (
	"os"
	"strings"

	"github.com/coeeter/ntmux/internal/tmux"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:                "ntmux",
	Short:              "Yet another tmux wrapper",
	Args:               cobra.ArbitraryArgs,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		if isHelpCommand(args) {
			printUnifiedHelp(cmd)
			return
		}

		if hasNtmuxConfigFileInRoot() {
			ApplyCmd.Run(cmd, []string{})
			return
		}

		tmux.PassThrough(args)
	},
}

func isHelpCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}

	for _, arg := range args {
		if arg == "-h" || arg == "-help" || arg == "--help" || arg == "help" {
			return true
		}
	}
	return false
}

func printUnifiedHelp(cmd *cobra.Command) {
	cmd.Help()

	cmd.Println("\nTmux Help:")
	output, err := tmux.PassThroughWithOutput([]string{"-h"})
	if err != nil {
		return
	}

	outputStr := strings.ReplaceAll(string(output), "tmux", "ntmux")
	cmd.Println(outputStr)
}

func hasNtmuxConfigFileInRoot() bool {
	_, err := os.Stat("ntmux.json")
	if err == nil {
		return true
	}
	_, err = os.Stat("ntmux.yaml")
	if err == nil {
		return true
	}
	return false
}

func init() {
	RootCmd.AddCommand(ApplyCmd)
}
