//go:build unit

package toolchain_test

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/toolchain"
	"github.com/rios0rios0/langforge/test/builders"
	"github.com/rios0rios0/langforge/test/doubles"
)

// versionRunner returns a runner whose version command prints output and
// remembers the command line it was asked to run in ran.
func versionRunner(output string, ran *[]string) *doubles.RunnerStub {
	return &doubles.RunnerStub{
		RunOutputFunc: func(_, name string, args ...string) (string, error) {
			*ran = append(append(*ran, name), args...)
			return output, nil
		},
	}
}

func TestRuntimeManager_Description(t *testing.T) {
	t.Parallel()

	t.Run("should describe the SDK it was given when asked", func(t *testing.T) {
		t.Parallel()

		// given
		sdk := builders.NewSDKBuilder().
			WithName("Tool").
			WithVersionManager("tvm").
			WithStartCommand("tool start").
			WithStopCommand("tool stop").
			Build().(toolchain.SDK)
		m := toolchain.NewRuntimeManager(&doubles.RunnerStub{}, sdk)

		// when
		name, manager, start, stop := m.SDKName(), m.VersionManager(), m.StartCommand(), m.StopCommand()

		// then
		assert.Equal(t, "Tool", name)
		assert.Equal(t, "tvm", manager)
		assert.Equal(t, "tool start", start)
		assert.Equal(t, "tool stop", stop)
	})

	t.Run("should report no stop command when the SDK has none", func(t *testing.T) {
		t.Parallel()

		// given
		sdk := builders.NewSDKBuilder().Build().(toolchain.SDK)
		m := toolchain.NewRuntimeManager(&doubles.RunnerStub{}, sdk)

		// when
		stop := m.StopCommand()

		// then
		assert.Empty(t, stop)
	})

	t.Run("should place the version in the install command when asked", func(t *testing.T) {
		t.Parallel()

		// given
		sdk := builders.NewSDKBuilder().WithInstallCommand("tvm install tool-%s").Build().(toolchain.SDK)
		m := toolchain.NewRuntimeManager(&doubles.RunnerStub{}, sdk)

		// when
		cmd := m.InstallCommand("1.2.3")

		// then
		assert.Equal(t, "tvm install tool-1.2.3", cmd)
	})
}

func TestRuntimeManager_CurrentVersion(t *testing.T) {
	t.Parallel()

	t.Run("should capture the version from the version command when the SDK is installed", func(t *testing.T) {
		t.Parallel()

		// given
		var ran []string
		sdk := builders.NewSDKBuilder().WithVersionCommand("tool", "--version").Build().(toolchain.SDK)
		m := toolchain.NewRuntimeManager(versionRunner("tool 1.2.3 (stable)", &ran), sdk)

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Equal(t, "1.2.3", version)
		assert.Equal(t, []string{"tool", "--version"}, ran)
	})

	t.Run("should return an empty version when the SDK binary is not found", func(t *testing.T) {
		t.Parallel()

		// given
		sdk := builders.NewSDKBuilder().Build().(toolchain.SDK)
		m := toolchain.NewRuntimeManager(doubles.NewRunnerStubBinaryNotFound(), sdk)

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Empty(t, version)
	})

	t.Run("should propagate the error when the version command fails for another reason", func(t *testing.T) {
		t.Parallel()

		// given
		expectedErr := errors.New("permission denied")
		runner := &doubles.RunnerStub{
			RunOutputFunc: func(_, _ string, _ ...string) (string, error) {
				return "", expectedErr
			},
		}
		sdk := builders.NewSDKBuilder().Build().(toolchain.SDK)
		m := toolchain.NewRuntimeManager(runner, sdk)

		// when
		version, err := m.CurrentVersion()

		// then
		require.ErrorIs(t, err, expectedErr)
		assert.NotErrorIs(t, err, exec.ErrNotFound)
		assert.Empty(t, version)
	})

	t.Run("should return an empty version when the output does not match the pattern", func(t *testing.T) {
		t.Parallel()

		// given
		var ran []string
		sdk := builders.NewSDKBuilder().Build().(toolchain.SDK)
		m := toolchain.NewRuntimeManager(versionRunner("unexpected output", &ran), sdk)

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Empty(t, version)
	})
}
