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
