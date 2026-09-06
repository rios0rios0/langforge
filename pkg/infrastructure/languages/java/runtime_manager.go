// Package java holds what the Java/Gradle and Java/Maven providers share: the
// JDK. The two ecosystems differ in their build tool, not in their SDK, its
// version manager, or the command that reports its version, so they compose the
// runtime manager defined here instead of each carrying an identical copy.
package java

import (
	"regexp"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// versionRe matches the version in `java --version` output, whose first line
// reads `openjdk 21.0.2 2024-01-16`.
var versionRe = regexp.MustCompile(`(\d+\.\d+(?:\.\d+)?)`)

// RuntimeManager provides SDK and runtime information for JVM projects. The
// command that starts an application is the only part the build tool decides,
// so it is the only part the caller supplies.
type RuntimeManager struct {
	*toolchain.RuntimeManager
}

// NewRuntimeManager creates a RuntimeManager that runs commands through runner
// and reports startCommand as the way to start a project.
func NewRuntimeManager(runner cmdexec.Runner, startCommand string) *RuntimeManager {
	return &RuntimeManager{toolchain.NewRuntimeManager(runner, jdk(startCommand))}
}

// jdk describes the Java SDK: sdkman installs it, the build tool's startCommand
// starts a project, there is no standard stop command, and java reports the
// version.
func jdk(startCommand string) toolchain.SDK {
	return toolchain.SDK{
		Name:           "Java",
		VersionManager: "sdkman",
		InstallCommand: "sdk install java %s",
		StartCommand:   startCommand,
		VersionCommand: cmdexec.CommandSpec{Name: "java", Args: []string{"--version"}},
		VersionPattern: versionRe,
	}
}
