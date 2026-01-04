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

		for _, session := range templ.Sessions {
			firstWindow := session.Windows[0]
			tmux.NewSession(session.Name, session.Dir, firstWindow.Name, firstWindow.Cmd, true)
			for i, window := range session.Windows {
				if i == 0 {
					continue
				}
				tmux.NewWindow(session.Name, window.Name, window.Dir, window.Cmd)
			}
			defaultWindow := session.Windows[0]
			for _, window := range session.Windows {
				if window.Default {
					defaultWindow = window
					break
				}
			}
			tmux.SelectWindow(session.Name, defaultWindow.Name)
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

func getTemplatePath(args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	if hasNtmuxConfigFileInRoot() {
		file, err := os.Stat("ntmux.json")
		if err == nil && !file.IsDir() {
			return "ntmux.json", nil
		}
		file, err = os.Stat("ntmux.yaml")
		if err == nil && !file.IsDir() {
			return "ntmux.yaml", nil
		}
	}

	return "", os.ErrNotExist
}
