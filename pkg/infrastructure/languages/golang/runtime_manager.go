package golang

import (
	"fmt"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// RuntimeManager provides SDK and runtime information for Go projects.
type RuntimeManager struct {
	runner cmdexec.Runner
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{runner: runner}
}

// SDKName returns the human-readable SDK name.
func (m *RuntimeManager) SDKName() string { return "Go" }

// VersionManager returns the version manager name.
func (m *RuntimeManager) VersionManager() string { return "gvm" }

// StartCommand returns the default command to run a Go project.
func (m *RuntimeManager) StartCommand() string { return "go run ." }

// StopCommand returns an empty string since Go has no long-running dev server.
func (m *RuntimeManager) StopCommand() string { return "" }

// InstallCommand returns the gvm command to install a specific Go version.
func (m *RuntimeManager) InstallCommand(version string) string {
	return fmt.Sprintf("gvm install go%s", version)
}

// CurrentVersion returns the currently installed Go version, or empty if not installed.
func (m *RuntimeManager) CurrentVersion() (string, error) {
	return cmdexec.CapturedVersion(m.runner, regexp.MustCompile(`go(\d+\.\d+(?:\.\d+)?)`), "go", "version")
}
