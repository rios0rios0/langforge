//go:build unit

package dart_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/dart"
	"github.com/rios0rios0/langforge/test/doubles"
)

// recordingRunner captures the commands a component would have executed.
type recordingRunner struct {
	*doubles.RunnerStub
	commands []string
}

func newRecordingRunner() *recordingRunner {
	r := &recordingRunner{}
	r.RunnerStub = &doubles.RunnerStub{
		RunFunc: func(_, name string, args ...string) error {
			r.commands = append(r.commands, strings.Join(append([]string{name}, args...), " "))
			return nil
		},
	}
	return r
}

func TestDartDependencyUpdater(t *testing.T) {
	t.Parallel()

	t.Run("should run dart pub for a plain Dart package", func(t *testing.T) {
		t.Parallel()

		// given
		dir, _ := writePubspec(t, "name: cli\n\ndependencies:\n  args: ^2.7.0\n")
		runner := newRecordingRunner()
		u := dart.NewDependencyUpdater(runner)

		// when
		err := u.UpdateAll(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{
			"dart pub upgrade --major-versions",
			"dart pub get",
		}, runner.commands)
	})

	t.Run("should run flutter pub for a Flutter project", func(t *testing.T) {
		t.Parallel()

		// given
		dir, _ := writePubspec(t, commentedPubspec)
		runner := newRecordingRunner()
		u := dart.NewDependencyUpdater(runner)

		// when
		err := u.UpdateAll(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{
			"flutter pub upgrade --major-versions",
			"flutter pub get",
		}, runner.commands)
	})

	t.Run("should list the manifest and the lockfile as files changed", func(t *testing.T) {
		t.Parallel()

		// given
		dir, path := writePubspec(t, commentedPubspec)
		lock := filepath.Join(dir, "pubspec.lock")
		require.NoError(t, os.WriteFile(lock, []byte("packages: {}\n"), 0o600))
		u := dart.NewDependencyUpdater(newRecordingRunner())

		// when
		files, err := u.FilesChanged(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{path, lock}, files)
	})

	t.Run("should omit the lockfile when the project has none", func(t *testing.T) {
		t.Parallel()

		// given
		dir, path := writePubspec(t, commentedPubspec)
		u := dart.NewDependencyUpdater(newRecordingRunner())

		// when
		files, err := u.FilesChanged(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{path}, files)
	})

	t.Run("should describe the upgrade commands for logging", func(t *testing.T) {
		t.Parallel()

		// given
		u := dart.NewDependencyUpdater(newRecordingRunner())

		// when
		commands := u.Commands()

		// then
		assert.Equal(t, []string{
			"dart pub upgrade --major-versions",
			"dart pub get",
		}, commands)
	})
}

func TestDartBuildValidator(t *testing.T) {
	t.Parallel()

	t.Run("should resolve then analyse with the Flutter toolchain", func(t *testing.T) {
		t.Parallel()

		// given
		dir, _ := writePubspec(t, commentedPubspec)
		runner := newRecordingRunner()
		v := dart.NewBuildValidator(runner)

		// when
		err := v.Validate(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{"flutter pub get", "flutter analyze"}, runner.commands)
	})

	t.Run("should resolve then analyse with the Dart toolchain", func(t *testing.T) {
		t.Parallel()

		// given
		dir, _ := writePubspec(t, "name: cli\n\ndependencies:\n  args: ^2.7.0\n")
		runner := newRecordingRunner()
		v := dart.NewBuildValidator(runner)

		// when
		err := v.Validate(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{"dart pub get", "dart analyze"}, runner.commands)
	})
}

func TestDartRuntimeManager(t *testing.T) {
	t.Parallel()

	t.Run("should parse the version from dart --version output", func(t *testing.T) {
		t.Parallel()

		// given
		runner := &doubles.RunnerStub{
			RunOutputFunc: func(_, _ string, _ ...string) (string, error) {
				return `Dart SDK version: 3.13.0 (stable) (Wed Aug 5 00:28:05 2026 -0700) on "linux_x64"`, nil
			},
		}
		m := dart.NewRuntimeManager(runner)

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Equal(t, "3.13.0", version)
	})

	t.Run("should return an empty version when the SDK is not installed", func(t *testing.T) {
		t.Parallel()

		// given
		m := dart.NewRuntimeManager(doubles.NewRunnerStubBinaryNotFound())

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Empty(t, version)
	})

	t.Run("should describe the SDK and its version manager", func(t *testing.T) {
		t.Parallel()

		// given
		m := dart.NewRuntimeManager(newRecordingRunner())

		// when
		name, manager, install := m.SDKName(), m.VersionManager(), m.InstallCommand("3.13.0")

		// then
		assert.Equal(t, "Dart", name)
		assert.Equal(t, "fvm", manager)
		assert.Equal(t, "fvm install 3.13.0", install)
	})
}
