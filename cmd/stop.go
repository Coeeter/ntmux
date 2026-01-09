package cmd

import (
	"os"

	"github.com/coeeter/ntmux/internal/template"
	"github.com/coeeter/ntmux/internal/tmux"
	"github.com/spf13/cobra"
)

var StopCmd = &cobra.Command{
	Use:   "stop [template-file]",
	Short: "Stop a tmux session template",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := getTemplatePath(args)
		if err != nil {
			cmd.Println("Error: No template file specified and no ntmux.json or ntmux.yaml found in the current directory.")
			return
		}

		cwd, err := os.Getwd()
		if err != nil {
			cmd.Println("Error getting current working directory:", err)
			return
		}

		templ, err := template.LoadTemplateFromFile(path, cwd)
		if err != nil {
			cmd.Println("Error loading template:", err)
			return
		}

		runner := tmux.NewRunner(tmux.GetShell())
		for _, session := range templ.Sessions {
			if tmux.HasSession(session.Name) {
				runner.KillSession(session.Name)
			}
		}
		runner.Execute()
	},
}
