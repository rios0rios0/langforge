package csharp

import (
	"regexp"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// dotnetVersionRe matches the version in `dotnet --version` output, which is the
// bare version on a line of its own, `8.0.404`: it is kept whole, so a preview
// suffix such as `9.0.100-preview.7.24407.12` survives.
var dotnetVersionRe = regexp.MustCompile(`^(\S+)`)

// RuntimeManager provides SDK and runtime information for C# projects.
type RuntimeManager struct {
	*toolchain.RuntimeManager
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{toolchain.NewRuntimeManager(runner, dotnetSDK())}
}

// dotnetSDK describes the .NET SDK: the dotnet CLI installs it, `dotnet run`
// starts a project, there is no standard stop command, and dotnet reports the
// version.
func dotnetSDK() toolchain.SDK {
	return toolchain.SDK{
		Name:           ".NET",
		VersionManager: dotnetCLI,
		InstallCommand: "dotnet sdk install %s",
		StartCommand:   "dotnet run",
		VersionCommand: cmdexec.CommandSpec{Name: dotnetCLI, Args: []string{"--version"}},
		VersionPattern: dotnetVersionRe,
	}
}
