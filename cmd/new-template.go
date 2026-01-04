package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/coeeter/ntmux/internal/template"
	"github.com/invopop/yaml"
	"github.com/spf13/cobra"
)

var format string

var NewTemplateCmd = &cobra.Command{
	Use:   "new-template",
	Short: "Create a new ntmux template file",
	Long: `Create a new ntmux.json or ntmux.yaml template file in the current directory.
If a custom template exists at ~/.config/ntmux/template.json or ~/.config/ntmux/template.yaml,
it will be used as the base. Otherwise, a default template will be created.`,
	Run: func(cmd *cobra.Command, args []string) {
		format = strings.ToLower(format)
		if format != "json" && format != "yaml" {
			cmd.Println("Error: format must be 'json' or 'yaml'")
			return
		}

		outputFile := fmt.Sprintf("ntmux.%s", format)
		if _, err := os.Stat(outputFile); err == nil {
			cmd.Printf("Error: %s already exists in the current directory.\n", outputFile)
			return
		}

		templ := loadCustomTemplate()
		if templ == nil {
			cwd, err := os.Getwd()
			if err != nil {
				cmd.Printf("Error getting current working directory: %v\n", err)
				return
			}
			templ = getDefaultTemplate(filepath.Base(cwd))
		}

		if err := writeTemplate(templ, outputFile); err != nil {
			cmd.Printf("Error writing template: %v\n", err)
			return
		}

		cmd.Printf("Created %s successfully!\n", outputFile)
	},
}

func init() {
	NewTemplateCmd.Flags().StringVarP(&format, "format", "f", "json", "Format of the template file (json or yaml)")
}

func loadCustomTemplate() *template.Template {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	configDir := filepath.Join(homeDir, ".config", "ntmux")

	jsonPath := filepath.Join(configDir, "template.json")
	if data, err := os.ReadFile(jsonPath); err == nil {
		var templ template.Template
		if err := json.Unmarshal(data, &templ); err == nil {
			return &templ
		}
	}

	yamlPath := filepath.Join(configDir, "template.yaml")
	if data, err := os.ReadFile(yamlPath); err == nil {
		var templ template.Template
		if err := yaml.Unmarshal(data, &templ); err == nil {
			return &templ
		}
	}

	ymlPath := filepath.Join(configDir, "template.yml")
	if data, err := os.ReadFile(ymlPath); err == nil {
		var templ template.Template
		if err := yaml.Unmarshal(data, &templ); err == nil {
			return &templ
		}
	}

	return nil
}

func getDefaultTemplate(dir string) *template.Template {
	return &template.Template{
		Schema: "https://raw.githubusercontent.com/coeeter/ntmux/main/schema.json",
		Sessions: []template.Session{
			{
				Name:    dir,
				Default: true,
				Windows: []template.Window{
					{
						Name:    "editor",
						Cmd:     "nvim .",
						Default: true,
					},
					{
						Name: "terminal",
					},
				},
			},
		},
	}
}

func writeTemplate(templ *template.Template, outputFile string) error {
	var data []byte
	var err error

	if strings.HasSuffix(outputFile, ".yaml") || strings.HasSuffix(outputFile, ".yml") {
		templ.Schema = "" // Remove schema for YAML files as yaml doesn't support $schema
		data, err = yaml.Marshal(templ)
	} else {
		data, err = json.MarshalIndent(templ, "", "  ")
	}

	if err != nil {
		return err
	}

	return os.WriteFile(outputFile, data, 0644)
}
