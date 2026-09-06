package toolchain

import (
	"fmt"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// SDK describes a language's software development kit: what it is called, which
// version manager installs it, how a project is started and stopped, and how
// the installed version is read.
type SDK struct {
	// Name is the human-readable SDK name, such as "Go" or "Node.js".
	Name string
	// VersionManager is the tool that installs SDK versions, such as "gvm" or "nvm".
	VersionManager string
	// InstallCommand is the version manager's install command with a single %s
	// placeholder that receives the version, such as "nvm install %s".
	InstallCommand string
	// StartCommand is the default command that runs a project.
	StartCommand string
	// StopCommand is the default command that stops a project, left empty when
	// the ecosystem has no standard one.
	StopCommand string
	// VersionCommand is the command whose output reports the installed version.
	VersionCommand cmdexec.CommandSpec
	// VersionPattern captures the version from that output in its first group.
	VersionPattern *regexp.Regexp
}

// RuntimeManager provides SDK and runtime information for a language from the
// SDK description it is given, reading the installed version through a Runner.
type RuntimeManager struct {
	runner cmdexec.Runner
	sdk    SDK
}

// NewRuntimeManager creates a RuntimeManager for sdk that runs commands through runner.
func NewRuntimeManager(runner cmdexec.Runner, sdk SDK) *RuntimeManager {
	return &RuntimeManager{runner: runner, sdk: sdk}
}

// SDKName returns the human-readable SDK name.
func (m *RuntimeManager) SDKName() string { return m.sdk.Name }

// VersionManager returns the version manager name.
func (m *RuntimeManager) VersionManager() string { return m.sdk.VersionManager }

// StartCommand returns the default command to run a project.
func (m *RuntimeManager) StartCommand() string { return m.sdk.StartCommand }

// StopCommand returns the default command to stop a project, or an empty string
// when the ecosystem has none.
func (m *RuntimeManager) StopCommand() string { return m.sdk.StopCommand }

// InstallCommand returns the version manager's command to install version.
func (m *RuntimeManager) InstallCommand(version string) string {
	return fmt.Sprintf(m.sdk.InstallCommand, version)
}

// CurrentVersion returns the installed SDK version, or an empty string when the
// SDK is not installed or reports a version the pattern cannot read.
func (m *RuntimeManager) CurrentVersion() (string, error) {
	command := m.sdk.VersionCommand
	return cmdexec.CapturedVersion(m.runner, m.sdk.VersionPattern, command.Name, command.Args...)
}
