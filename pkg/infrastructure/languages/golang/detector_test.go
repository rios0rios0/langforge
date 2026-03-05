package golang_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/golang"
)

func TestGoDetector(t *testing.T) {
	t.Parallel()

	t.Run("should detect Go project when go.mod exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/foo\n\ngo 1.21\n"), 0o600))
		d := &golang.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.True(t, detected)
	})

	t.Run("should not detect Go project when go.mod is absent", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		d := &golang.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.False(t, detected)
	})

	t.Run("should return LanguageGo", func(t *testing.T) {
		t.Parallel()

		d := &golang.Detector{}
		assert.Equal(t, entities.LanguageGo, d.Language())
	})
}
