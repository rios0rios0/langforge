//go:build unit

package java_test

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/java"
	"github.com/rios0rios0/langforge/test/doubles"
)

func TestRuntimeManager_CurrentVersion(t *testing.T) {
	t.Parallel()

	t.Run("should return parsed version when java --version succeeds", func(t *testing.T) {
		t.Parallel()

		// given
		runner := &doubles.RunnerStub{
			RunOutputFunc: func(_, _ string, _ ...string) (string, error) {
				return "openjdk 21.0.2 2024-01-16", nil
			},
		}
		m := java.NewRuntimeManager(runner, "mvn spring-boot:run")

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Equal(t, "21.0.2", version)
	})

	t.Run("should return empty string when binary is not found", func(t *testing.T) {
		t.Parallel()

		// given
		runner := doubles.NewRunnerStubBinaryNotFound()
		m := java.NewRuntimeManager(runner, "mvn spring-boot:run")

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
		m := java.NewRuntimeManager(runner, "mvn spring-boot:run")

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
		m := java.NewRuntimeManager(runner, "mvn spring-boot:run")

		// when
		version, err := m.CurrentVersion()

		// then
		require.NoError(t, err)
		assert.Empty(t, version)
	})
}

func TestRuntimeManager_StartCommand(t *testing.T) {
	t.Parallel()

	t.Run("should report the build tool's start command", func(t *testing.T) {
		t.Parallel()

		// given
		gradle := java.NewRuntimeManager(&doubles.RunnerStub{}, "./gradlew bootRun")
		maven := java.NewRuntimeManager(&doubles.RunnerStub{}, "mvn spring-boot:run")

		// when
		gradleStart, mavenStart := gradle.StartCommand(), maven.StartCommand()

		// then
		assert.Equal(t, "./gradlew bootRun", gradleStart)
		assert.Equal(t, "mvn spring-boot:run", mavenStart)
	})

	t.Run("should report the same SDK for either build tool", func(t *testing.T) {
		t.Parallel()

		// given
		gradle := java.NewRuntimeManager(&doubles.RunnerStub{}, "./gradlew bootRun")
		maven := java.NewRuntimeManager(&doubles.RunnerStub{}, "mvn spring-boot:run")

		// when & then
		assert.Equal(t, "Java", gradle.SDKName())
		assert.Equal(t, "Java", maven.SDKName())
		assert.Equal(t, "sdkman", gradle.VersionManager())
		assert.Equal(t, "sdkman", maven.VersionManager())
		assert.Equal(t, "sdk install java 21.0.2", maven.InstallCommand("21.0.2"))
		assert.Empty(t, maven.StopCommand())
	})
}
