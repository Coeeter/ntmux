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

func SendKeys(target, keys string) {
	PassThrough([]string{"send-keys", "-t", target, keys, "C-m"})
}

func NewSession(sessionName, rootDir, windowName, firstWindowCmd string, detached bool) {
	args := []string{"new-session", "-s", sessionName, "-c", rootDir}
	if detached {
		args = append(args, "-d")
	}
	if windowName != "" {
		args = append(args, "-n", windowName)
	}
	PassThrough(args)
	if firstWindowCmd != "" {
		SendKeys(sessionName+":"+windowName, firstWindowCmd)
	}
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
		SendKeys(sessionName+":"+windowName, command)
	}
}

func SelectWindow(sessionName, windowName string) {
	PassThrough([]string{"select-window", "-t", sessionName + ":" + windowName})
}
