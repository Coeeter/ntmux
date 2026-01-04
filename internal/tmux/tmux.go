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

func HasSession(sessionName string) bool {
	_, err := PassThroughWithOutput([]string{"has-session", "-t", sessionName})
	return err == nil
}

func NewSession(sessionName, rootDir, windowName, firstWindowCmd string, detached bool) {
	args := []string{"new-session", "-s", sessionName, "-c", rootDir}
	if detached {
		args = append(args, "-d")
	}
	if windowName != "" {
		args = append(args, "-n", windowName)
	}
	if firstWindowCmd != "" {
		shell := GetShell()
		args = append(args, shell.GetCompleteCmd(firstWindowCmd))
	}

	PassThrough(args)
}

func NewWindow(sessionName, windowName, rootDir, command string) {
	args := []string{"new-window", "-t", sessionName, "-n", windowName, "-c", rootDir}

	if command != "" {
		shell := GetShell()
		args = append(args, shell.GetCompleteCmd(command))
	}

	PassThrough(args)
}

func SelectWindow(sessionName, windowName string) {
	PassThrough([]string{"select-window", "-t", sessionName + ":" + windowName})
}

func AttachSession(sessionName string) {
	if IsInTmux() {
		PassThrough([]string{"switch-client", "-t", sessionName})
	} else {
		PassThrough([]string{"attach-session", "-t", sessionName})
	}
}
