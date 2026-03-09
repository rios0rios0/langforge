//go:build unit

package dockerfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/dockerfile"
)

func TestDockerfileDetector(t *testing.T) {
	t.Parallel()

	t.Run("should detect project when Dockerfile exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte("FROM alpine:3.19\n"), 0o600))
		d := &dockerfile.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.True(t, detected)
	})

	t.Run("should detect project when Dockerfile variant exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "Dockerfile.dev"), []byte("FROM golang:1.26\n"), 0o600))
		d := &dockerfile.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.True(t, detected)
	})

	t.Run("should not detect project when no Dockerfile exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		d := &dockerfile.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.False(t, detected)
	})

	t.Run("should return LanguageDockerfile", func(t *testing.T) {
		t.Parallel()

		// given
		d := &dockerfile.Detector{}

		// then
		assert.Equal(t, entities.LanguageDockerfile, d.Language())
	})
}
