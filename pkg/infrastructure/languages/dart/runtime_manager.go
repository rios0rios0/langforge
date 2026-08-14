package dart

import (
	"fmt"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

const minVersionMatchGroups = 2

// dartVersionRe matches the version in `dart --version` output, which reads
// `Dart SDK version: 3.13.0 (stable) (...) on "linux_x64"` and goes to stdout.
var dartVersionRe = regexp.MustCompile(`Dart SDK version:\s+(\d+\.\d+\.\d+)`)

// RuntimeManager provides SDK and runtime information for Dart projects.
type RuntimeManager struct {
	runner cmdexec.Runner
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{runner: runner}
}

// SDKName returns the human-readable SDK name.
func (m *RuntimeManager) SDKName() string { return "Dart" }

// VersionManager returns the version manager name.
//
// fvm manages Flutter installations, and every Flutter install carries the Dart
// SDK it was built against, so it is the version manager for both. Pure Dart has
// no widely adopted equivalent.
func (m *RuntimeManager) VersionManager() string { return "fvm" }

// StartCommand returns the default command to run a Dart project.
// A Flutter application is started with `flutter run` instead.
func (m *RuntimeManager) StartCommand() string { return "dart run" }

// StopCommand returns an empty string since there is no standard stop command.
func (m *RuntimeManager) StopCommand() string { return "" }

// InstallCommand returns the fvm command to install a specific SDK version.
func (m *RuntimeManager) InstallCommand(version string) string {
	return fmt.Sprintf("fvm install %s", version)
}

// CurrentVersion returns the installed Dart SDK version, or empty if absent.
func (m *RuntimeManager) CurrentVersion() (string, error) {
	output, err := m.runner.RunOutput(".", "dart", "--version")
	if err != nil {
		if cmdexec.IsBinaryNotFound(err) {
			return "", nil
		}
		return "", err
	}
	matches := dartVersionRe.FindStringSubmatch(output)
	if len(matches) < minVersionMatchGroups {
		return "", nil
	}
	return matches[1], nil
}
