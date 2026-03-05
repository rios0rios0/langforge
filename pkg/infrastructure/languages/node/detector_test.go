package node_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/node"
)

func TestNodeDetector(t *testing.T) {
	t.Parallel()

	t.Run("should detect Node project when package.json exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"name":"test","version":"1.0.0"}`), 0o600))
		d := &node.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.True(t, detected)
	})

	t.Run("should not detect Node project when package.json is absent", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		d := &node.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.False(t, detected)
	})

	t.Run("should return LanguageNode", func(t *testing.T) {
		t.Parallel()

		d := &node.Detector{}
		assert.Equal(t, entities.LanguageNode, d.Language())
	})
}
