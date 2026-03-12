//go:build unit

package fileutil_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

func TestIsGlobPattern(t *testing.T) {
	t.Parallel()

	t.Run("should return true for asterisk patterns", func(t *testing.T) {
		t.Parallel()

		// given / when / then
		assert.True(t, fileutil.IsGlobPattern("*.tf"))
		assert.True(t, fileutil.IsGlobPattern("src/*.go"))
	})

	t.Run("should return true for question mark patterns", func(t *testing.T) {
		t.Parallel()

		// given / when / then
		assert.True(t, fileutil.IsGlobPattern("file?.txt"))
	})

	t.Run("should return true for bracket patterns", func(t *testing.T) {
		t.Parallel()

		// given / when / then
		assert.True(t, fileutil.IsGlobPattern("[abc].go"))
	})

	t.Run("should return false for exact paths", func(t *testing.T) {
		t.Parallel()

		// given / when / then
		assert.False(t, fileutil.IsGlobPattern("go.mod"))
		assert.False(t, fileutil.IsGlobPattern("package.json"))
		assert.False(t, fileutil.IsGlobPattern("pyproject.toml"))
	})
}

func TestExtractExtension(t *testing.T) {
	t.Parallel()

	t.Run("should extract extension from glob pattern", func(t *testing.T) {
		t.Parallel()

		// given / when / then
		assert.Equal(t, ".tf", fileutil.ExtractExtension("*.tf"))
		assert.Equal(t, ".hcl", fileutil.ExtractExtension("*.hcl"))
		assert.Equal(t, ".go", fileutil.ExtractExtension("src/*.go"))
	})

	t.Run("should extract extension from base name when directory has dots", func(t *testing.T) {
		t.Parallel()

		// given / when / then
		assert.Equal(t, ".go", fileutil.ExtractExtension("dir.with.dot/*.go"))
	})

	t.Run("should return input when directory has dots but base has no extension", func(t *testing.T) {
		t.Parallel()

		// given / when / then
		assert.Equal(t, "dir.with.dot/*", fileutil.ExtractExtension("dir.with.dot/*"))
	})

	t.Run("should return input when no dot is found", func(t *testing.T) {
		t.Parallel()

		// given / when / then
		assert.Equal(t, "Makefile", fileutil.ExtractExtension("Makefile"))
	})
}

func TestNewFileChecker(t *testing.T) {
	t.Parallel()

	t.Run("should call exactCheck for exact paths", func(t *testing.T) {
		t.Parallel()

		// given
		exactCalled := false
		checker := fileutil.NewFileChecker(
			func(path string) (bool, error) {
				exactCalled = true
				assert.Equal(t, "go.mod", path)
				return true, nil
			},
			func(pattern string) (bool, error) {
				t.Fatal("globCheck should not be called for exact paths")
				return false, nil
			},
		)

		// when
		found, err := checker("go.mod")

		// then
		require.NoError(t, err)
		assert.True(t, found)
		assert.True(t, exactCalled)
	})

	t.Run("should call globCheck for glob patterns", func(t *testing.T) {
		t.Parallel()

		// given
		globCalled := false
		checker := fileutil.NewFileChecker(
			func(path string) (bool, error) {
				t.Fatal("exactCheck should not be called for glob patterns")
				return false, nil
			},
			func(pattern string) (bool, error) {
				globCalled = true
				assert.Equal(t, "*.tf", pattern)
				return true, nil
			},
		)

		// when
		found, err := checker("*.tf")

		// then
		require.NoError(t, err)
		assert.True(t, found)
		assert.True(t, globCalled)
	})

	t.Run("should propagate errors from exactCheck", func(t *testing.T) {
		t.Parallel()

		// given
		expectedErr := errors.New("API error")
		checker := fileutil.NewFileChecker(
			func(path string) (bool, error) { return false, expectedErr },
			func(pattern string) (bool, error) { return false, nil },
		)

		// when
		_, err := checker("go.mod")

		// then
		require.ErrorIs(t, err, expectedErr)
	})

	t.Run("should propagate errors from globCheck", func(t *testing.T) {
		t.Parallel()

		// given
		expectedErr := errors.New("API error")
		checker := fileutil.NewFileChecker(
			func(path string) (bool, error) { return false, nil },
			func(pattern string) (bool, error) { return false, expectedErr },
		)

		// when
		_, err := checker("*.tf")

		// then
		require.ErrorIs(t, err, expectedErr)
	})
}
