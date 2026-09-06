package golang

import (
	"regexp"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// goVersionRe matches the version in `go version` output, which reads
// `go version go1.23.4 linux/amd64`.
var goVersionRe = regexp.MustCompile(`go(\d+\.\d+(?:\.\d+)?)`)

// RuntimeManager provides SDK and runtime information for Go projects.
type RuntimeManager struct {
	*toolchain.RuntimeManager
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{toolchain.NewRuntimeManager(runner, goSDK())}
}

// goSDK describes the Go SDK: gvm installs it, `go run` starts a project, there
// is no long-running dev server to stop, and `go version` reports the version.
func goSDK() toolchain.SDK {
	return toolchain.SDK{
		Name:           "Go",
		VersionManager: "gvm",
		InstallCommand: "gvm install go%s",
		StartCommand:   "go run .",
		VersionCommand: cmdexec.CommandSpec{Name: "go", Args: []string{"version"}},
		VersionPattern: goVersionRe,
	}
}
