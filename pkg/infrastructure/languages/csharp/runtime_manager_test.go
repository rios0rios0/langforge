//go:build unit

package csharp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/csharp"
	"github.com/rios0rios0/langforge/test/doubles"
)

func TestRuntimeManager_CurrentVersion(t *testing.T) {
	t.Parallel()

	t.Run("should return the reported version when dotnet --version succeeds", func(t *testing.T) {
		t.Parallel()

		// given
		var ran []string
		runner := &doubles.RunnerStub{
			RunOutputFunc: func(_, name string, args ...string) (string, error) {
				ran = append(append(ran, name), args...)
				return "8.0.404", nil
			},
		}
		m := csharp.NewRuntimeManager(runner)

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Equal(t, "8.0.404", version)
		assert.Equal(t, []string{"dotnet", "--version"}, ran)
	})

	t.Run("should keep the preview suffix when a preview SDK is installed", func(t *testing.T) {
		t.Parallel()

		// given
		runner := &doubles.RunnerStub{
			RunOutputFunc: func(_, _ string, _ ...string) (string, error) {
				return "9.0.100-preview.7.24407.12", nil
			},
		}
		m := csharp.NewRuntimeManager(runner)

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Equal(t, "9.0.100-preview.7.24407.12", version)
	})

	t.Run("should return empty string when binary is not found", func(t *testing.T) {
		t.Parallel()

		// given
		m := csharp.NewRuntimeManager(doubles.NewRunnerStubBinaryNotFound())

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Empty(t, version)
	})
}

func TestRuntimeManager_Description(t *testing.T) {
	t.Parallel()

	t.Run("should describe the .NET SDK and its dotnet commands", func(t *testing.T) {
		t.Parallel()

		// given
		m := csharp.NewRuntimeManager(&doubles.RunnerStub{})

		// when & then
		assert.Equal(t, ".NET", m.SDKName())
		assert.Equal(t, "dotnet", m.VersionManager())
		assert.Equal(t, "dotnet sdk install 8.0.404", m.InstallCommand("8.0.404"))
		assert.Equal(t, "dotnet run", m.StartCommand())
		assert.Empty(t, m.StopCommand())
	})
}
