package cmdexec

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
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

// IsBinaryNotFound reports whether the error indicates a missing executable.
func IsBinaryNotFound(err error) bool {
	return errors.Is(err, exec.ErrNotFound)
}

// minVersionMatchGroups is the smallest submatch slice that carries a captured
// version: the whole match plus capture group 1.
const minVersionMatchGroups = 2

// CapturedVersion runs a version command and returns the first capture group of
// pattern, which every RuntimeManager uses to read its SDK version.
//
// An absent binary is not an error: it means the SDK is not installed, which the
// RuntimeManager contract reports as an empty version. Output that the pattern
// does not match is treated the same way, since a version this code cannot parse
// is no more usable than none at all.
func CapturedVersion(runner Runner, pattern *regexp.Regexp, name string, args ...string) (string, error) {
	output, err := runner.RunOutput(".", name, args...)
	if err != nil {
		if IsBinaryNotFound(err) {
			return "", nil
		}
		return "", err
	}

	matches := pattern.FindStringSubmatch(output)
	if len(matches) < minVersionMatchGroups {
		return "", nil
	}
	return matches[1], nil
}

// CommandSpec defines a command with its name and arguments.
type CommandSpec struct {
	Name string
	Args []string
}

// String returns the command as a single displayable string.
func (c CommandSpec) String() string {
	return strings.Join(append([]string{c.Name}, c.Args...), " ")
}

// CommandStrings converts a slice of CommandSpec to displayable strings.
func CommandStrings(specs []CommandSpec) []string {
	out := make([]string, len(specs))
	for i, s := range specs {
		out[i] = s.String()
	}
	return out
}
