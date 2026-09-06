package dart

import (
	"regexp"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
)

// dartVersionRe matches the version in `dart --version` output, which reads
// `Dart SDK version: 3.13.0 (stable) (...) on "linux_x64"` and goes to stdout.
var dartVersionRe = regexp.MustCompile(`Dart SDK version:\s+(\d+\.\d+\.\d+)`)

// RuntimeManager provides SDK and runtime information for Dart projects.
type RuntimeManager struct {
	*toolchain.RuntimeManager
}

// NewRuntimeManager creates a RuntimeManager with the given runner.
func NewRuntimeManager(runner cmdexec.Runner) *RuntimeManager {
	return &RuntimeManager{toolchain.NewRuntimeManager(runner, dartSDK())}
}

// dartSDK describes the Dart SDK.
//
// fvm is its version manager: fvm manages Flutter installations, and every
// Flutter install carries the Dart SDK it was built against, so it is the
// version manager for both. Pure Dart has no widely adopted equivalent.
//
// `dart run` is the default start command; a Flutter application is started
// with `flutter run` instead. There is no standard stop command.
func dartSDK() toolchain.SDK {
	return toolchain.SDK{
		Name:           "Dart",
		VersionManager: "fvm",
		InstallCommand: "fvm install %s",
		StartCommand:   "dart run",
		VersionCommand: cmdexec.CommandSpec{Name: "dart", Args: []string{"--version"}},
		VersionPattern: dartVersionRe,
	}
}
