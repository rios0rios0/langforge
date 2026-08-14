//go:build unit

package dart_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/dart"
)

func TestBumpBuildNumber(t *testing.T) {
	t.Parallel()

	t.Run("should return the new version when there is no build number", func(t *testing.T) {
		t.Parallel()

		// given
		current := "0.1.0"

		// when
		result := dart.BumpBuildNumber(current, "0.2.0")

		// then
		assert.Equal(t, "0.2.0", result)
	})

	t.Run("should increment the build number when one is present", func(t *testing.T) {
		t.Parallel()

		// given
		current := "1.0.0+1"

		// when
		result := dart.BumpBuildNumber(current, "1.1.0")

		// then
		assert.Equal(t, "1.1.0+2", result)
	})

	t.Run("should preserve zero padding when incrementing the build number", func(t *testing.T) {
		t.Parallel()

		// given
		current := "2.10.2+021002"

		// when
		result := dart.BumpBuildNumber(current, "2.11.0")

		// then
		assert.Equal(t, "2.11.0+021003", result)
	})

	t.Run("should widen the build number when the increment outgrows its padding", func(t *testing.T) {
		t.Parallel()

		// given
		current := "1.0.0+99"

		// when
		result := dart.BumpBuildNumber(current, "1.0.1")

		// then
		assert.Equal(t, "1.0.1+100", result)
	})

	t.Run("should carry a non-numeric build suffix through unchanged", func(t *testing.T) {
		t.Parallel()

		// given
		current := "1.0.0+nightly"

		// when
		result := dart.BumpBuildNumber(current, "1.1.0")

		// then
		assert.Equal(t, "1.1.0+nightly", result)
	})
}

func TestReplaceVersion(t *testing.T) {
	t.Parallel()

	t.Run("should replace the version and keep the trailing comment", func(t *testing.T) {
		t.Parallel()

		// given
		content := "name: app\nversion: 2.10.2+021002 # See README.md for details on versioning.\n"

		// when
		result, ok := dart.ReplaceVersion(content, "2.11.0")

		// then
		assert.True(t, ok)
		assert.Equal(
			t,
			"name: app\nversion: 2.11.0+021003 # See README.md for details on versioning.\n",
			result,
		)
	})

	t.Run("should leave the SDK constraint and dependencies untouched", func(t *testing.T) {
		t.Parallel()

		// given
		content := "name: medhub\nversion: 0.1.0\n\n" +
			"environment:\n  sdk: ^3.13.0\n\n" +
			"dependencies:\n  go_router: ^17.0.0\n  intl: 0.20.2\n"

		// when
		result, ok := dart.ReplaceVersion(content, "0.2.0")

		// then
		assert.True(t, ok)
		assert.Contains(t, result, "version: 0.2.0\n")
		assert.Contains(t, result, "  sdk: ^3.13.0\n")
		assert.Contains(t, result, "  go_router: ^17.0.0\n")
		assert.Contains(t, result, "  intl: 0.20.2\n")
	})

	t.Run("should report no match when the manifest declares no version", func(t *testing.T) {
		t.Parallel()

		// given
		content := "name: workspace_member\nenvironment:\n  sdk: ^3.13.0\n"

		// when
		result, ok := dart.ReplaceVersion(content, "1.0.0")

		// then
		assert.False(t, ok)
		assert.Equal(t, content, result)
	})

	t.Run("should replace only the first version line", func(t *testing.T) {
		t.Parallel()

		// given
		content := "version: 1.0.0\nname: app\nversion: 9.9.9\n"

		// when
		result, ok := dart.ReplaceVersion(content, "1.1.0")

		// then
		assert.True(t, ok)
		assert.Equal(t, "version: 1.1.0\nname: app\nversion: 9.9.9\n", result)
	})
}
