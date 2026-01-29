package cmd

import (
	"os"
	"strings"

	"github.com/coeeter/ntmux/internal/template"
	"github.com/coeeter/ntmux/internal/tmux"
	"github.com/spf13/cobra"
)

var sessionsFlag string

var ApplyCmd = &cobra.Command{
	Use:   "apply [template-file]",
	Short: "Apply a tmux session template",
	Run: func(cmd *cobra.Command, args []string) {
		path, err := getTemplatePath(args)
		if err != nil {
			cmd.Println("Error: No template file specified and no ntmux.json, ntmux.yaml, or ntmux.yml found in the current directory.")
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

		// Parse session filter if provided
		sessionFilter := make(map[string]bool)
		if sessionsFlag != "" {
			for _, name := range strings.Split(sessionsFlag, ",") {
				sessionFilter[strings.TrimSpace(name)] = true
			}
		}

		// Validate that all specified sessions exist in the template
		if len(sessionFilter) > 0 {
			templateSessions := make(map[string]bool)
			for _, session := range templ.Sessions {
				templateSessions[session.Name] = true
			}
			var invalidSessions []string
			for name := range sessionFilter {
				if !templateSessions[name] {
					invalidSessions = append(invalidSessions, name)
				}
			}
			if len(invalidSessions) > 0 {
				cmd.Printf("Error: session(s) not found in template: %s\n", strings.Join(invalidSessions, ", "))
				cmd.Println("Available sessions:", strings.Join(getSessionNames(templ.Sessions), ", "))
				return
			}
		}

		for _, session := range templ.Sessions {
			// Skip if session filter is set and this session is not in it
			if len(sessionFilter) > 0 && !sessionFilter[session.Name] {
				continue
			}
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
			// Only consider sessions that match the filter
			if len(sessionFilter) > 0 && !sessionFilter[session.Name] {
				continue
			}
			if session.Default {
				defaultSession = session.Name
				break
			}
		}

		if defaultSession != "" {
			if tmux.IsInTmux() {
				runner.SwitchClient(defaultSession)
			} else {
				runner.AttachSession(defaultSession)
			}
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
	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == "ntmux.yml" {
			return "ntmux.yml", nil
		}
	}

	return "", os.ErrNotExist
}

func init() {
	ApplyCmd.Flags().StringVarP(&sessionsFlag, "sessions", "s", "", "Comma-separated list of session names to create (e.g., --sessions=frontend,backend)")
}

func getSessionNames(sessions []template.Session) []string {
	names := make([]string, len(sessions))
	for i, s := range sessions {
		names[i] = s.Name
	}
	return names
}
