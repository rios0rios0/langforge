package javamaven

import (
	"fmt"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

const minVersionMatchGroups = 2

// RuntimeManager provides SDK and runtime information for Java/Maven projects.
type RuntimeManager struct {
	runner cmdexec.Runner
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{runner: runner}
}

// SDKName returns the human-readable SDK name.
func (m *RuntimeManager) SDKName() string { return "Java" }

// VersionManager returns the version manager name.
func (m *RuntimeManager) VersionManager() string { return "sdkman" }

// StartCommand returns the default command to run a Java/Maven project.
func (m *RuntimeManager) StartCommand() string { return "mvn spring-boot:run" }

// StopCommand returns an empty string since there is no standard stop command.
func (m *RuntimeManager) StopCommand() string { return "" }

// InstallCommand returns the sdkman command to install a specific Java version.
func (m *RuntimeManager) InstallCommand(version string) string {
	return fmt.Sprintf("sdk install java %s", version)
}

// CurrentVersion returns the currently installed Java version, or empty if not installed.
func (m *RuntimeManager) CurrentVersion() (string, error) {
	output, err := m.runner.RunOutput(".", "java", "--version")
	if err != nil {
		if cmdexec.IsBinaryNotFound(err) {
			return "", nil
		}
		return "", err
	}
	re := regexp.MustCompile(`(\d+\.\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < minVersionMatchGroups {
		return "", nil
	}
	return matches[1], nil
}
