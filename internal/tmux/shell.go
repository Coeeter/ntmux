package tmux

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type Shell struct {
	BinPath string
	Name    string
	os      string
}

func GetShell() *Shell {
	shellPath := os.Getenv("SHELL")
	if shellPath == "" {
		shellPath = "/bin/sh"
	}
	shellName := extractShellName(shellPath)
	return &Shell{
		BinPath: shellPath,
		Name:    shellName,
		os:      detectOS(),
	}
}

func extractShellName(shellPath string) string {
	parts := strings.Split(shellPath, "/")
	return parts[len(parts)-1]
}

func detectOS() string {
	return runtime.GOOS
}

func (s *Shell) GetCompleteCmd(cmd string) string {
	cmd = strings.TrimSpace(cmd)

	if strings.ToLower(s.os) == "windows" {
		cmd = strings.ReplaceAll(cmd, `"`, `\"`)
		if s.Name == "powershell" || s.Name == "pwsh" {
			return fmt.Sprintf("%s -NoExit -Command \"& {%s}\"", s.BinPath, cmd)
		} else {
			return fmt.Sprintf("%s /K \"%s\"", s.BinPath, cmd)
		}
	}

	// Unix
	cmd = strings.ReplaceAll(cmd, `'`, `'\''`)
	return fmt.Sprintf("%s -c '%s; exec %s'", s.BinPath, cmd, s.BinPath)
}
