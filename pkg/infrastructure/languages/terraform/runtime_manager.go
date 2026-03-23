package terraform

import (
	"fmt"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// RuntimeManager provides SDK and runtime information for Terraform projects.
type RuntimeManager struct {
	runner cmdexec.Runner
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{runner: runner}
}

// SDKName returns the human-readable SDK name.
func (m *RuntimeManager) SDKName() string { return "Terraform" }

// VersionManager returns the version manager name.
func (m *RuntimeManager) VersionManager() string { return "tfenv" }

// StartCommand returns the default command to apply a Terraform project.
func (m *RuntimeManager) StartCommand() string { return "terraform apply" }

// StopCommand returns an empty string since there is no standard stop command.
func (m *RuntimeManager) StopCommand() string { return "" }

// InstallCommand returns the tfenv command to install a specific Terraform version.
func (m *RuntimeManager) InstallCommand(version string) string {
	return fmt.Sprintf("tfenv install %s", version)
}

// CurrentVersion returns the currently installed Terraform version, or empty if not installed.
func (m *RuntimeManager) CurrentVersion() (string, error) {
	output, err := m.runner.RunOutput(".", "terraform", "version")
	if err != nil {
		return "", nil
	}
	re := regexp.MustCompile(`Terraform\s+v(\d+\.\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", nil
	}
	return matches[1], nil
}
