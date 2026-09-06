package ruby

import (
	"regexp"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// rubyVersionRe matches the version in `ruby --version` output, which reads
// `ruby 3.3.0 (2023-12-25 revision 5124f9ac75) [x86_64-linux]`.
var rubyVersionRe = regexp.MustCompile(`ruby\s+(\d+\.\d+(?:\.\d+)?)`)

// RuntimeManager provides SDK and runtime information for Ruby projects.
type RuntimeManager struct {
	*toolchain.RuntimeManager
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{toolchain.NewRuntimeManager(runner, rubySDK())}
}

// rubySDK describes the Ruby SDK: rbenv installs it, the Rails server starts a
// project, there is no standard stop command, and ruby reports the version.
func rubySDK() toolchain.SDK {
	return toolchain.SDK{
		Name:           "Ruby",
		VersionManager: "rbenv",
		InstallCommand: "rbenv install %s",
		StartCommand:   "bundle exec rails server",
		VersionCommand: cmdexec.CommandSpec{Name: "ruby", Args: []string{"--version"}},
		VersionPattern: rubyVersionRe,
	}
}
