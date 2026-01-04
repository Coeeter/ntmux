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

func HasSession(sessionName string) bool {
	output, err := PassThroughWithOutput([]string{"has-session", "-t", sessionName})
	if err != nil {
		return false
	}
	return string(output) == ""
}

func IsInTmux() bool {
	return os.Getenv("TMUX") != ""
}

type TmuxRunner struct {
	shell    *Shell
	commands [][]string
}

func NewRunner(shell *Shell) *TmuxRunner {
	return &TmuxRunner{
		shell:    shell,
		commands: make([][]string, 0),
	}
}

func (r *TmuxRunner) NewSession(sessionName, rootDir, windowName, firstWindowCmd string, detached bool) {
	args := []string{"new-session", "-s", sessionName, "-c", rootDir}
	if detached {
		args = append(args, "-d")
	}
	if windowName != "" {
		args = append(args, "-n", windowName)
	}
	if firstWindowCmd != "" {
		args = append(args, r.shell.GetCompleteCmd(firstWindowCmd))
	}
	r.commands = append(r.commands, args)
}

func (r *TmuxRunner) NewWindow(sessionName, windowName, rootDir, command string) {
	args := []string{"new-window", "-t", sessionName, "-n", windowName, "-c", rootDir}
	if command != "" {
		args = append(args, r.shell.GetCompleteCmd(command))
	}
	r.commands = append(r.commands, args)
}

func (r *TmuxRunner) SelectWindow(sessionName, windowName string) {
	args := []string{"select-window", "-t", sessionName + ":" + windowName}
	r.commands = append(r.commands, args)
}

func (r *TmuxRunner) AttachSession(sessionName string) {
	args := []string{"attach-session", "-t", sessionName}
	r.commands = append(r.commands, args)
}

func (r *TmuxRunner) Execute() {
	if len(r.commands) == 0 {
		return
	}

	var args []string
	for i, cmd := range r.commands {
		if i > 0 {
			args = append(args, ";")
		}
		args = append(args, cmd...)
	}

	PassThrough(args)
}
