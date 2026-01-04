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

		shell := tmux.GetShell()
		runner := tmux.NewRunner(shell)

		for _, session := range templ.Sessions {
			if tmux.HasSession(session.Name) {
				continue
			}
			firstWindow := session.Windows[0]
			runner.NewSession(session.Name, session.Dir, firstWindow.Name, firstWindow.Cmd, true)

			for i, window := range session.Windows {
				if i == 0 {
					continue
				}
				runner.NewWindow(session.Name, window.Name, window.Dir, window.Cmd)
			}

			defaultWindow := session.Windows[0]
			for _, window := range session.Windows {
				if window.Default {
					defaultWindow = window
					break
				}
			}
			runner.SelectWindow(session.Name, defaultWindow.Name)
		}

		var defaultSession string
		for _, session := range templ.Sessions {
			if session.Default {
				defaultSession = session.Name
				break
			}
		}

		if defaultSession != "" {
			runner.AttachSession(defaultSession)
		}

		runner.Execute()
	},
}

func getTemplatePath(args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	entries, err := os.ReadDir(".")
	if err != nil {
		return "", os.ErrNotExist
	}

	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == "ntmux.json" {
			return "ntmux.json", nil
		}
	}
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == "ntmux.yaml" {
			return "ntmux.yaml", nil
		}
	}

	return "", os.ErrNotExist
}
