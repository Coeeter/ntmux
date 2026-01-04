package template

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/invopop/yaml"
)

type Window struct {
	Name    string `json:"name" jsonSchema:"title=Window Name,description=Name of the tmux window"`
	Dir     string `json:"dir,omitempty" jsonSchema:"title=Directory,description=Working directory for the window"`
	Cmd     string `json:"cmd,omitempty" jsonSchema:"title=Command,description=Command to run in the window"`
	Default bool   `json:"default,omitempty" jsonSchema:"title=Default Window,description=Whether this window should be the default one in the session"`
}

type Session struct {
	Name    string   `json:"name" jsonSchema:"title=Session Name,description=Name of the tmux session"`
	Dir     string   `json:"dir,omitempty" jsonSchema:"title=Directory,description=Working directory for the session"`
	Windows []Window `json:"windows" jsonSchema:"title=Windows,description=List of windows in the session"`
	Default bool     `json:"default,omitempty" jsonSchema:"title=Default Session,description=Whether this session should be the default one to attach to"`
}

type Template struct {
	Schema   string    `json:"$schema,omitempty" jsonSchema:"title=JSON Schema,description=The JSON schema version"`
	Sessions []Session `json:"sessions" jsonSchema:"title=Sessions,description=List of tmux sessions to create"`
}

func LoadTemplateFromFile(path string, cwd string) (*Template, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var template Template
	if isYAML(path) {
		err = yaml.Unmarshal(buffer, &template)
	} else {
		err = json.Unmarshal(buffer, &template)
	}
	if err != nil {
		return nil, err
	}

	for i, session := range template.Sessions {
		if session.Dir == "" {
			template.Sessions[i].Dir = cwd
		} else {
			template.Sessions[i].Dir = filepath.Join(cwd, session.Dir)
		}

		for j, window := range session.Windows {
			if window.Dir == "" {
				template.Sessions[i].Windows[j].Dir = template.Sessions[i].Dir
			} else {
				template.Sessions[i].Windows[j].Dir = filepath.Join(cwd, window.Dir)
			}
		}
	}

	var hasDefault bool
	for _, session := range template.Sessions {
		if session.Default {
			hasDefault = true
			break
		}
	}

	if !hasDefault && len(template.Sessions) > 0 {
		template.Sessions[0].Default = true
	}

	return &template, nil
}

func isYAML(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".yaml" || ext == ".yml"
}
