package golang_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/golang"
)

func TestGoVersionReader(t *testing.T) {
	t.Parallel()

	t.Run("should read version from go.mod version comment", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		content := "// version: 1.2.3\nmodule example.com/foo\n\ngo 1.21\n"
		require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0o600))
		r := &golang.VersionReader{}

		// when
		version, err := r.ReadVersion(dir)

		// then
		require.NoError(t, err)
		assert.Equal(t, "1.2.3", version.String())
	})

	t.Run("should return error when no version comment in go.mod", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		content := "module example.com/foo\n\ngo 1.21\n"
		require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte(content), 0o600))
		r := &golang.VersionReader{}

		// when
		_, err := r.ReadVersion(dir)

		// then
		require.Error(t, err)
	})
}
