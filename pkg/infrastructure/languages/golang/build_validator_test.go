//go:build unit

package golang_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/golang"
	"github.com/rios0rios0/langforge/test/doubles"
)

func TestBuildValidator_Validate(t *testing.T) {
	t.Parallel()

	t.Run("should run lint then build commands in order", func(t *testing.T) {
		t.Parallel()

		// given
		var calledCommands []string
		runner := &doubles.RunnerStub{
			RunFunc: func(_, name string, _ ...string) error {
				calledCommands = append(calledCommands, name)
				return nil
			},
		}
		v := golang.NewBuildValidator(runner)

		// when
		err := v.Validate("/repo")

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{"golangci-lint", "go"}, calledCommands)
	})

	t.Run("should return error when lint command fails", func(t *testing.T) {
		t.Parallel()

		// given
		lintErr := errors.New("lint failed")
		runner := &doubles.RunnerStub{
			RunFunc: func(_, _ string, _ ...string) error {
				return lintErr
			},
		}
		v := golang.NewBuildValidator(runner)

		// when
		err := v.Validate("/repo")

		// then
		require.ErrorIs(t, err, lintErr)
	})

	t.Run("should return error when build command fails", func(t *testing.T) {
		t.Parallel()

		// given
		buildErr := errors.New("build failed")
		callCount := 0
		runner := &doubles.RunnerStub{
			RunFunc: func(_, _ string, _ ...string) error {
				callCount++
				if callCount == 2 {
					return buildErr
				}
				return nil
			},
		}
		v := golang.NewBuildValidator(runner)

		// when
		err := v.Validate("/repo")

		// then
		require.ErrorIs(t, err, buildErr)
	})
}

func TestBuildValidator_Commands(t *testing.T) {
	t.Parallel()

	t.Run("should return lint commands as displayable strings", func(t *testing.T) {
		t.Parallel()

		// given
		v := golang.NewBuildValidator(&doubles.RunnerStub{})

		// when
		cmds := v.LintCommands()

		// then
		assert.Equal(t, []string{"golangci-lint run ./..."}, cmds)
	})

	t.Run("should return build commands as displayable strings", func(t *testing.T) {
		t.Parallel()

		// given
		v := golang.NewBuildValidator(&doubles.RunnerStub{})

		// when
		cmds := v.BuildCommands()

		// then
		assert.Equal(t, []string{"go build ./..."}, cmds)
	})
}
