package exec

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Runner executes shell commands.
type Runner interface {
	// Run executes a command in the given working directory.
	Run(dir string, name string, args ...string) error
}

// DefaultRunner is the default shell command runner using os/exec.
type DefaultRunner struct{}

// NewDefaultRunner creates a new DefaultRunner.
func NewDefaultRunner() *DefaultRunner {
	return &DefaultRunner{}
}

// Run executes the command with the given args in the specified directory.
func (r *DefaultRunner) Run(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...) // #nosec G204
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command %q failed: %w\nstderr: %s", name, err, stderr.String())
	}
	return nil
}
