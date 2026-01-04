package cmd

import (
	"os"

	"github.com/coeeter/ntmux/internal/template"
	"github.com/coeeter/ntmux/internal/tmux"
	"github.com/spf13/cobra"
)

var ApplyCmd = &cobra.Command{
	Use:   "apply [template-file]",
	Short: "Apply a tmux session template",
	Run: func(cmd *cobra.Command, args []string) {
		path := "ntmux.yaml"
		if len(args) > 0 {
			path = args[0]
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

		for _, session := range templ.Sessions {
			tmux.NewSession(session.Name, session.Dir, true)
			for _, window := range session.Windows {
				tmux.NewWindow(session.Name, window.Name, window.Dir, window.Cmd)
			}
		}

		var defaultSession string
		for _, session := range templ.Sessions {
			if session.Default {
				defaultSession = session.Name
				break
			}
		}

		if defaultSession != "" {
			tmux.AttachSession(defaultSession)
		}
	},
}
