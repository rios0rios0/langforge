package terraform

import (
	"fmt"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// terraformCLIVersionRe matches the version the installed binary reports, whose
// first line reads `Terraform v1.7.0`. It is deliberately distinct from
// terraformVersionRe, which reads the required_version constraint a repository
// declares.
var terraformCLIVersionRe = regexp.MustCompile(`Terraform\s+v(\d+\.\d+(?:\.\d+)?)`)

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
	return cmdexec.CapturedVersion(m.runner, terraformCLIVersionRe, "terraform", "version")
}
