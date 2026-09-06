//go:build unit

package toolchain_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/pkg/support/cmdexec"
	"github.com/rios0rios0/langforge/test/doubles"
)

// lintThenBuild describes a toolchain with two lint checks and one build step.
func lintThenBuild() toolchain.Commands {
	return toolchain.Commands{
		Lint: []cmdexec.CommandSpec{
			{Name: "fmt", Args: []string{"-check"}},
			{Name: "lint", Args: []string{"./..."}},
		},
		Build: []cmdexec.CommandSpec{{Name: "build", Args: []string{"-o", "out"}}},
	}
}

// recordingRunner returns a runner that appends every command line it is asked
// to run to ran, and fails the command named failing with failErr.
func recordingRunner(ran *[]string, failing string, failErr error) *doubles.RunnerStub {
	return &doubles.RunnerStub{
		RunFunc: func(_, name string, args ...string) error {
			*ran = append(*ran, strings.Join(append([]string{name}, args...), " "))
			if name == failing {
				return failErr
			}
			return nil
		},
	}
}

func TestBuildValidator_Validate(t *testing.T) {
	t.Parallel()

	t.Run("should run every lint command before any build command when validating", func(t *testing.T) {
		t.Parallel()

		// given
		var ran []string
		v := toolchain.NewBuildValidator(recordingRunner(&ran, "", nil), lintThenBuild())

		// when
		err := v.Validate("/repo")

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{"fmt -check", "lint ./...", "build -o out"}, ran)
	})

	t.Run("should run every command in the repository path when validating", func(t *testing.T) {
		t.Parallel()

		// given
		var dirs []string
		runner := &doubles.RunnerStub{
			RunFunc: func(dir, _ string, _ ...string) error {
				dirs = append(dirs, dir)
				return nil
			},
		}
		v := toolchain.NewBuildValidator(runner, lintThenBuild())

		// when
		err := v.Validate("/repo")

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{"/repo", "/repo", "/repo"}, dirs)
	})

	t.Run("should stop at the failing lint command when one fails", func(t *testing.T) {
		t.Parallel()

		// given
		var ran []string
		lintErr := errors.New("lint failed")
		v := toolchain.NewBuildValidator(recordingRunner(&ran, "lint", lintErr), lintThenBuild())

		// when
		err := v.Validate("/repo")

		// then
		require.ErrorIs(t, err, lintErr)
		assert.Equal(t, []string{"fmt -check", "lint ./..."}, ran)
	})

	t.Run("should return the build error when a build command fails", func(t *testing.T) {
		t.Parallel()

		// given
		var ran []string
		buildErr := errors.New("build failed")
		v := toolchain.NewBuildValidator(recordingRunner(&ran, "build", buildErr), lintThenBuild())

		// when
		err := v.Validate("/repo")

		// then
		require.ErrorIs(t, err, buildErr)
		assert.Equal(t, []string{"fmt -check", "lint ./...", "build -o out"}, ran)
	})

	t.Run("should succeed without running anything when no commands are configured", func(t *testing.T) {
		t.Parallel()

		// given
		var ran []string
		v := toolchain.NewBuildValidator(recordingRunner(&ran, "", nil), toolchain.Commands{})

		// when
		err := v.Validate("/repo")

		// then
		require.NoError(t, err)
		assert.Empty(t, ran)
	})
}

func TestBuildValidator_Commands(t *testing.T) {
	t.Parallel()

	t.Run("should return lint commands as displayable strings when asked", func(t *testing.T) {
		t.Parallel()

		// given
		v := toolchain.NewBuildValidator(&doubles.RunnerStub{}, lintThenBuild())

		// when
		cmds := v.LintCommands()

		// then
		assert.Equal(t, []string{"fmt -check", "lint ./..."}, cmds)
	})

	t.Run("should return build commands as displayable strings when asked", func(t *testing.T) {
		t.Parallel()

		// given
		v := toolchain.NewBuildValidator(&doubles.RunnerStub{}, lintThenBuild())

		// when
		cmds := v.BuildCommands()

		// then
		assert.Equal(t, []string{"build -o out"}, cmds)
	})

	t.Run("should return no commands when a phase has none", func(t *testing.T) {
		t.Parallel()

		// given
		lintOnly := toolchain.Commands{Lint: lintThenBuild().Lint}
		v := toolchain.NewBuildValidator(&doubles.RunnerStub{}, lintOnly)

		// when
		cmds := v.BuildCommands()

		// then
		assert.Empty(t, cmds)
	})
}
