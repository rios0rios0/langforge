package ruby

import (
	"fmt"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

const minVersionMatchGroups = 2

// RuntimeManager provides SDK and runtime information for Ruby projects.
type RuntimeManager struct {
	runner cmdexec.Runner
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{runner: runner}
}

// SDKName returns the human-readable SDK name.
func (m *RuntimeManager) SDKName() string { return "Ruby" }

// VersionManager returns the version manager name.
func (m *RuntimeManager) VersionManager() string { return "rbenv" }

// StartCommand returns the default command to start a Ruby project.
func (m *RuntimeManager) StartCommand() string { return "bundle exec rails server" }

// StopCommand returns an empty string since there is no standard stop command.
func (m *RuntimeManager) StopCommand() string { return "" }

// InstallCommand returns the rbenv command to install a specific Ruby version.
func (m *RuntimeManager) InstallCommand(version string) string {
	return fmt.Sprintf("rbenv install %s", version)
}

// CurrentVersion returns the currently installed Ruby version, or empty if not installed.
func (m *RuntimeManager) CurrentVersion() (string, error) {
	output, err := m.runner.RunOutput(".", "ruby", "--version")
	if err != nil {
		if cmdexec.IsBinaryNotFound(err) {
			return "", nil
		}
		return "", err
	}
	re := regexp.MustCompile(`ruby\s+(\d+\.\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < minVersionMatchGroups {
		return "", nil
	}
	return matches[1], nil
}
