package cmdexec

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Runner executes shell commands.
type Runner interface {
	// Run executes a command in the given working directory.
	Run(dir string, name string, args ...string) error

	// RunOutput executes a command and returns its stdout output.
	RunOutput(dir string, name string, args ...string) (string, error)
}

// DefaultRunner is the default shell command runner using os/exec.
type DefaultRunner struct{}

// NewDefaultRunner creates a new DefaultRunner.
func NewDefaultRunner() *DefaultRunner {
	return &DefaultRunner{}
}

// Run executes the command with the given args in the specified directory.
func (r *DefaultRunner) Run(dir string, name string, args ...string) error {
	cmd := exec.CommandContext(context.Background(), name, args...) // #nosec G204
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command %q failed: %w\nstderr: %s", name, err, stderr.String())
	}
	return nil
}

// RunOutput executes the command and returns its stdout as a trimmed string.
func (r *DefaultRunner) RunOutput(dir string, name string, args ...string) (string, error) {
	cmd := exec.CommandContext(context.Background(), name, args...) // #nosec G204
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command %q failed: %w\nstderr: %s", name, err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}
