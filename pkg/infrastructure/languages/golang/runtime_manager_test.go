//go:build unit

package golang_test

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/golang"
	"github.com/rios0rios0/langforge/test/doubles"
)

func TestRuntimeManager_CurrentVersion(t *testing.T) {
	t.Parallel()

	t.Run("should return parsed version when go version succeeds", func(t *testing.T) {
		t.Parallel()

		// given
		runner := &doubles.RunnerStub{
			RunOutputFunc: func(_, _ string, _ ...string) (string, error) {
				return "go version go1.23.4 linux/amd64", nil
			},
		}
		m := golang.NewRuntimeManager(runner)

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Equal(t, "1.23.4", version)
	})

	t.Run("should return empty string when binary is not found", func(t *testing.T) {
		t.Parallel()

		// given
		runner := doubles.NewRunnerStubBinaryNotFound()
		m := golang.NewRuntimeManager(runner)

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Empty(t, version)
	})

	t.Run("should propagate non-binary-not-found errors", func(t *testing.T) {
		t.Parallel()

		// given
		expectedErr := errors.New("permission denied")
		runner := &doubles.RunnerStub{
			RunOutputFunc: func(_, _ string, _ ...string) (string, error) {
				return "", expectedErr
			},
		}
		m := golang.NewRuntimeManager(runner)

		// when
		version, err := m.CurrentVersion()

		// then
		require.Error(t, err)
		assert.Empty(t, version)
		assert.NotErrorIs(t, err, exec.ErrNotFound)
	})

	t.Run("should return empty string when version output has no match", func(t *testing.T) {
		t.Parallel()

		// given
		runner := &doubles.RunnerStub{
			RunOutputFunc: func(_, _ string, _ ...string) (string, error) {
				return "unexpected output", nil
			},
		}
		m := golang.NewRuntimeManager(runner)

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Empty(t, version)
	})
}
