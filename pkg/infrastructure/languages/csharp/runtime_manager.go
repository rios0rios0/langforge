package csharp

import (
	"fmt"
	"strings"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// RuntimeManager provides SDK and runtime information for C# projects.
type RuntimeManager struct {
	runner cmdexec.Runner
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{runner: runner}
}

// SDKName returns the human-readable SDK name.
func (m *RuntimeManager) SDKName() string { return ".NET" }

// VersionManager returns the version manager name.
func (m *RuntimeManager) VersionManager() string { return "dotnet" }

// StartCommand returns the default command to run a C# project.
func (m *RuntimeManager) StartCommand() string { return "dotnet run" }

// StopCommand returns an empty string since there is no standard stop command.
func (m *RuntimeManager) StopCommand() string { return "" }

// InstallCommand returns the command to install a specific .NET SDK version.
func (m *RuntimeManager) InstallCommand(version string) string {
	return fmt.Sprintf("dotnet sdk install %s", version)
}

// CurrentVersion returns the currently installed .NET SDK version, or empty if not installed.
func (m *RuntimeManager) CurrentVersion() (string, error) {
	output, err := m.runner.RunOutput(".", "dotnet", "--version")
	if err != nil {
		if cmdexec.IsBinaryNotFound(err) {
			return "", nil
		}
		return "", err
	}
	version := strings.TrimSpace(output)
	if version == "" {
		return "", nil
	}
	return version, nil
}
