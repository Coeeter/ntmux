package tmux

import (
	"os"
	"os/exec"
)

func PassThrough(args []string) {
	cmd := exec.Command("tmux", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func PassThroughWithOutput(args []string) ([]byte, error) {
	cmd := exec.Command("tmux", args...)
	return cmd.CombinedOutput()
}

func IsInTmux() bool {
	return os.Getenv("TMUX") != ""
}

func NewSession(sessionName, rootDir string, detached bool) {
	args := []string{"new-session", "-s", sessionName, "-c", rootDir}
	if detached {
		args = append(args, "-d")
	}
	PassThrough(args)
}

func AttachSession(sessionName string) {
	if IsInTmux() {
		PassThrough([]string{"switch-client", "-t", sessionName})
	} else {
		PassThrough([]string{"attach-session", "-t", sessionName})
	}
}

func NewWindow(sessionName, windowName, rootDir, command string) {
	args := []string{"new-window", "-t", sessionName, "-n", windowName, "-c", rootDir}
	PassThrough(args)

	if command != "" {
		args = []string{"send-keys", "-t", sessionName + ":" + windowName, command, "C-m"}
		PassThrough(args)
	}
}
