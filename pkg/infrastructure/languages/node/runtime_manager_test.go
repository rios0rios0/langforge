//go:build unit

package node_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/node"
	"github.com/rios0rios0/langforge/test/doubles"
)

func TestRuntimeManager_CurrentVersion(t *testing.T) {
	t.Parallel()

	t.Run("should drop the v prefix when node --version succeeds", func(t *testing.T) {
		t.Parallel()

		// given
		var ran []string
		runner := &doubles.RunnerStub{
			RunOutputFunc: func(_, name string, args ...string) (string, error) {
				ran = append(append(ran, name), args...)
				return "v20.11.0", nil
			},
		}
		m := node.NewRuntimeManager(runner)

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Equal(t, "20.11.0", version)
		assert.Equal(t, []string{"node", "--version"}, ran)
	})

	t.Run("should return empty string when binary is not found", func(t *testing.T) {
		t.Parallel()

		// given
		m := node.NewRuntimeManager(doubles.NewRunnerStubBinaryNotFound())

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Empty(t, version)
	})
}

func TestRuntimeManager_Description(t *testing.T) {
	t.Parallel()

	t.Run("should describe the Node.js SDK and its nvm and npm commands", func(t *testing.T) {
		t.Parallel()

		// given
		m := node.NewRuntimeManager(&doubles.RunnerStub{})

		// when & then
		assert.Equal(t, "Node.js", m.SDKName())
		assert.Equal(t, "nvm", m.VersionManager())
		assert.Equal(t, "nvm install 20.11.0", m.InstallCommand("20.11.0"))
		assert.Equal(t, "npm start", m.StartCommand())
		assert.Empty(t, m.StopCommand())
	})
}
