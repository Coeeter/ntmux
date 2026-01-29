package cmd

import (
	"os"
	"path/filepath"
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

		if hasNtmuxConfigFileInRoot() && len(args) == 0 {
			ApplyCmd.Run(cmd, []string{})
			return
		}

		// If no arguments are provided, create a new tmux session
		if len(args) == 0 {
			cwd, err := os.Getwd()
			if err != nil {
				tmux.PassThrough(args)
				return
			}

			dirName := filepath.Base(cwd)

			instanceExists := tmux.HasSession(dirName)
			runner := tmux.NewRunner(tmux.GetShell())
			if !instanceExists {
				runner.NewSession(dirName, cwd, "", "", true)
			}

			if tmux.IsInTmux() {
				runner.SwitchClient(dirName)
			} else {
				runner.AttachSession(dirName)
			}

			runner.Execute()
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
	entries, err := os.ReadDir(".")
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() && (entry.Name() == "ntmux.json" || entry.Name() == "ntmux.yaml" || entry.Name() == "ntmux.yml") {
			return true
		}
	}
	return false
}

func init() {
	RootCmd.AddCommand(ApplyCmd)
	RootCmd.AddCommand(NewTemplateCmd)
	RootCmd.AddCommand(StopCmd)
	RootCmd.AddCommand(PickCmd)
}
