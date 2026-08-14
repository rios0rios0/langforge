package javagradle

import (
	"fmt"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// RuntimeManager provides SDK and runtime information for Java/Gradle projects.
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

// StartCommand returns the default command to run a Java/Gradle project.
func (m *RuntimeManager) StartCommand() string { return "./gradlew bootRun" }

// StopCommand returns an empty string since there is no standard stop command.
func (m *RuntimeManager) StopCommand() string { return "" }

// InstallCommand returns the sdkman command to install a specific Java version.
func (m *RuntimeManager) InstallCommand(version string) string {
	return fmt.Sprintf("sdk install java %s", version)
}

// CurrentVersion returns the currently installed Java version, or empty if not installed.
func (m *RuntimeManager) CurrentVersion() (string, error) {
	return cmdexec.CapturedVersion(m.runner, regexp.MustCompile(`(\d+\.\d+(?:\.\d+)?)`), "java", "--version")
}
