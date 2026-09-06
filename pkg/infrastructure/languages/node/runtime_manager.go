package node

import (
	"regexp"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// nodeVersionRe matches the version in `node --version` output, which reads
// `v20.11.0`: the leading v is dropped and the rest is kept whole.
var nodeVersionRe = regexp.MustCompile(`^v?(\S+)`)

// RuntimeManager provides SDK and runtime information for Node.js projects.
type RuntimeManager struct {
	*toolchain.RuntimeManager
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{toolchain.NewRuntimeManager(runner, nodeSDK())}
}

// nodeSDK describes the Node.js SDK: nvm installs it, `npm start` starts a
// project, there is no standard stop command, and node reports the version.
func nodeSDK() toolchain.SDK {
	return toolchain.SDK{
		Name:           "Node.js",
		VersionManager: "nvm",
		InstallCommand: "nvm install %s",
		StartCommand:   "npm start",
		VersionCommand: cmdexec.CommandSpec{Name: "node", Args: []string{"--version"}},
		VersionPattern: nodeVersionRe,
	}
}
