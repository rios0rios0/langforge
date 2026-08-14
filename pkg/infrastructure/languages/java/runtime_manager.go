// Package java holds what the Java/Gradle and Java/Maven providers share: the
// JDK. The two ecosystems differ in their build tool, not in their SDK, its
// version manager, or the command that reports its version, so they compose the
// runtime manager defined here instead of each carrying an identical copy.
package java

import (
	"fmt"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// versionRe matches the version in `java --version` output, whose first line
// reads `openjdk 21.0.2 2024-01-16`.
var versionRe = regexp.MustCompile(`(\d+\.\d+(?:\.\d+)?)`)

// RuntimeManager provides SDK and runtime information for JVM projects. The
// command that starts an application is the only part the build tool decides,
// so it is the only part the caller supplies.
type RuntimeManager struct {
	runner       cmdexec.Runner
	startCommand string
}

// NewRuntimeManager creates a RuntimeManager that runs commands through runner
// and reports startCommand as the way to start a project.
func NewRuntimeManager(runner cmdexec.Runner, startCommand string) *RuntimeManager {
	return &RuntimeManager{runner: runner, startCommand: startCommand}
}

// SDKName returns the human-readable SDK name.
func (m *RuntimeManager) SDKName() string { return "Java" }

// VersionManager returns the version manager name.
func (m *RuntimeManager) VersionManager() string { return "sdkman" }

// StartCommand returns the build tool's command to run a project.
func (m *RuntimeManager) StartCommand() string { return m.startCommand }

// StopCommand returns an empty string since there is no standard stop command.
func (m *RuntimeManager) StopCommand() string { return "" }

// InstallCommand returns the sdkman command to install a specific Java version.
func (m *RuntimeManager) InstallCommand(version string) string {
	return fmt.Sprintf("sdk install java %s", version)
}

// CurrentVersion returns the currently installed Java version, or empty if not installed.
func (m *RuntimeManager) CurrentVersion() (string, error) {
	return cmdexec.CapturedVersion(m.runner, versionRe, "java", "--version")
}
