package python

import (
	"regexp"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// pythonVersionRe matches the version in `python3 --version` output, which
// reads `Python 3.12.1`.
var pythonVersionRe = regexp.MustCompile(`Python\s+(\d+\.\d+(?:\.\d+)?)`)

// RuntimeManager provides SDK and runtime information for Python projects.
type RuntimeManager struct {
	*toolchain.RuntimeManager
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{toolchain.NewRuntimeManager(runner, pythonSDK())}
}

// pythonSDK describes the Python SDK: pyenv installs it, `pdm start` starts a
// project, there is no standard stop command, and python3 reports the version.
func pythonSDK() toolchain.SDK {
	return toolchain.SDK{
		Name:           "Python",
		VersionManager: "pyenv",
		InstallCommand: "pyenv install %s",
		StartCommand:   "pdm start",
		VersionCommand: cmdexec.CommandSpec{Name: "python3", Args: []string{"--version"}},
		VersionPattern: pythonVersionRe,
	}
}
