//go:build unit

package versions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLTSRelease(t *testing.T) {
	t.Parallel()

	t.Run("should return true when LTS is a non-empty string", func(t *testing.T) {
		t.Parallel()

		// given
		release := nodeRelease{Version: "v22.0.0", LTS: "Jod"}

		// when
		result := isLTSRelease(release)

		// then
		assert.True(t, result)
	})

	t.Run("should return false when LTS is an empty string", func(t *testing.T) {
		t.Parallel()

		// given
		release := nodeRelease{Version: "v23.0.0", LTS: ""}

		// when
		result := isLTSRelease(release)

		// then
		assert.False(t, result)
	})

	t.Run("should return false when LTS is false", func(t *testing.T) {
		t.Parallel()

		// given
		release := nodeRelease{Version: "v23.0.0", LTS: false}

		// when
		result := isLTSRelease(release)

		// then
		assert.False(t, result)
	})

	t.Run("should return true when LTS is true", func(t *testing.T) {
		t.Parallel()

		// given
		release := nodeRelease{Version: "v22.0.0", LTS: true}

		// when
		result := isLTSRelease(release)

		// then
		assert.True(t, result)
	})

	t.Run("should return false when LTS is an unexpected type", func(t *testing.T) {
		t.Parallel()

		// given
		release := nodeRelease{Version: "v22.0.0", LTS: 42}

		// when
		result := isLTSRelease(release)

		// then
		assert.False(t, result)
	})
}

func TestIsActiveEOL(t *testing.T) {
	t.Parallel()

	t.Run("should return true when EOL is false", func(t *testing.T) {
		t.Parallel()

		// given
		eol := false

		// when
		result := isActiveEOL(eol)

		// then
		assert.True(t, result)
	})

	t.Run("should return false when EOL is true", func(t *testing.T) {
		t.Parallel()

		// given
		eol := true

		// when
		result := isActiveEOL(eol)

		// then
		assert.False(t, result)
	})

	t.Run("should return true when EOL is a future date string", func(t *testing.T) {
		t.Parallel()

		// given
		eol := "2099-12-31"

		// when
		result := isActiveEOL(eol)

		// then
		assert.True(t, result)
	})

	t.Run("should return false when EOL is a past date string", func(t *testing.T) {
		t.Parallel()

		// given
		eol := "2020-01-01"

		// when
		result := isActiveEOL(eol)

		// then
		assert.False(t, result)
	})

	t.Run("should return false when EOL is an invalid date string", func(t *testing.T) {
		t.Parallel()

		// given
		eol := "not-a-date"

		// when
		result := isActiveEOL(eol)

		// then
		assert.False(t, result)
	})

	t.Run("should return false when EOL is an unexpected type", func(t *testing.T) {
		t.Parallel()

		// given
		eol := 42

		// when
		result := isActiveEOL(eol)

		// then
		assert.False(t, result)
	})
}
