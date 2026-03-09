//go:build unit

package pipeline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/pipeline"
)

func TestPipelineDetector(t *testing.T) {
	t.Parallel()

	t.Run("should detect pipeline project when GitHub Actions workflow exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		workflowDir := filepath.Join(dir, ".github", "workflows")
		require.NoError(t, os.MkdirAll(workflowDir, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(workflowDir, "ci.yaml"), []byte("name: CI\n"), 0o600))
		d := &pipeline.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.True(t, detected)
	})

	t.Run("should detect pipeline project when Azure DevOps pipeline exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "azure-pipelines.yml"), []byte("trigger:\n"), 0o600))
		d := &pipeline.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.True(t, detected)
	})

	t.Run("should detect pipeline project when nested Azure DevOps template exists", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		adoDir := filepath.Join(dir, "azure-devops", "templates")
		require.NoError(t, os.MkdirAll(adoDir, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(adoDir, "build.yaml"), []byte("steps:\n"), 0o600))
		d := &pipeline.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.True(t, detected)
	})

	t.Run("should not detect pipeline project when no CI files exist", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		d := &pipeline.Detector{}

		// when
		detected, err := d.Detect(dir)

		// then
		require.NoError(t, err)
		assert.False(t, detected)
	})

	t.Run("should return LanguagePipeline", func(t *testing.T) {
		t.Parallel()

		// given
		d := &pipeline.Detector{}

		// then
		assert.Equal(t, entities.LanguagePipeline, d.Language())
	})
}
