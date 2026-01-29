package cmd

import (
	"os/exec"
	"strings"

	"github.com/coeeter/ntmux/internal/tmux"
	"github.com/spf13/cobra"
)

var PickCmd = &cobra.Command{
	Use:   "pick",
	Short: "Interactively pick and attach to a tmux session using fzf",
	Run: func(cmd *cobra.Command, args []string) {
		doesFxfExist, err := exec.LookPath("fzf")
		if err != nil || doesFxfExist == "" {
			cmd.Println("fzf is not installed or not found in PATH")
			return
		}

		output, err := tmux.PassThroughWithOutput([]string{"list-sessions", "-F", "#S"})
		if err != nil {
			cmd.Println("Error listing tmux sessions:", err)
			return
		}

		outputString := string(output)
		outputReader := strings.NewReader(outputString)
		fzfCmd := exec.Command("fzf", "--prompt", "Select tmux session: ")
		fzfCmd.Stdin = outputReader
		outputBytes, err := fzfCmd.Output()
		if err != nil {
			cmd.Println("Error selecting tmux session:", err)
			return
		}

		selectedSession := strings.TrimSpace(string(outputBytes))
		if selectedSession == "" {
			cmd.Println("No session selected.")
			return
		}

		tmux.PassThrough([]string{"attach-session", "-t", selectedSession})
	},
}
