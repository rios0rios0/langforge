//go:build unit

package repositories_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/domain/repositories"
)

// versionWriterStub is a VersionWriter that reports a canned file list.
type versionWriterStub struct {
	files []string
	err   error
}

func (s *versionWriterStub) WriteVersion(_ string, _ entities.Version) error { return nil }

func (s *versionWriterStub) FilesChanged(_ string) ([]string, error) { return s.files, s.err }

// dependencyUpdaterStub is a DependencyUpdater that reports a canned file list.
type dependencyUpdaterStub struct {
	files []string
	err   error
}

func (s *dependencyUpdaterStub) UpdateAll(_ string) error { return nil }

func (s *dependencyUpdaterStub) FilesChanged(_ string) ([]string, error) { return s.files, s.err }

func (s *dependencyUpdaterStub) Commands() []string { return nil }

func newComposite(writer *versionWriterStub, updater *dependencyUpdaterStub) *repositories.CompositeProvider {
	return &repositories.CompositeProvider{VersionWriter: writer, DependencyUpdater: updater}
}

func TestCompositeProviderFilesChanged(t *testing.T) {
	t.Parallel()

	t.Run("should return the union with the writer's files first", func(t *testing.T) {
		t.Parallel()

		// given
		writer := &versionWriterStub{files: []string{"pubspec.yaml"}}
		updater := &dependencyUpdaterStub{files: []string{"pubspec.lock"}}
		provider := newComposite(writer, updater)

		// when
		files, err := provider.FilesChanged("/repo")

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{"pubspec.yaml", "pubspec.lock"}, files)
	})

	t.Run("should report a file both halves claim only once", func(t *testing.T) {
		t.Parallel()

		// given — every ecosystem where the version and the dependencies live in
		// the same manifest reports it twice, and the caller must not act on it twice.
		writer := &versionWriterStub{files: []string{"pubspec.yaml"}}
		updater := &dependencyUpdaterStub{files: []string{"pubspec.yaml", "pubspec.lock"}}
		provider := newComposite(writer, updater)

		// when
		files, err := provider.FilesChanged("/repo")

		// then
		require.NoError(t, err)
		assert.Equal(t, []string{"pubspec.yaml", "pubspec.lock"}, files)
	})

	t.Run("should return no files when neither half reports any", func(t *testing.T) {
		t.Parallel()

		// given
		provider := newComposite(&versionWriterStub{}, &dependencyUpdaterStub{})

		// when
		files, err := provider.FilesChanged("/repo")

		// then
		require.NoError(t, err)
		assert.Empty(t, files)
	})

	t.Run("should propagate the error when the version writer fails", func(t *testing.T) {
		t.Parallel()

		// given
		expectedErr := errors.New("cannot stat pubspec.yaml")
		writer := &versionWriterStub{err: expectedErr}
		updater := &dependencyUpdaterStub{files: []string{"pubspec.lock"}}
		provider := newComposite(writer, updater)

		// when
		files, err := provider.FilesChanged("/repo")

		// then
		require.ErrorIs(t, err, expectedErr)
		assert.Nil(t, files)
	})

	t.Run("should propagate the error when the dependency updater fails", func(t *testing.T) {
		t.Parallel()

		// given
		expectedErr := errors.New("cannot stat pubspec.lock")
		writer := &versionWriterStub{files: []string{"pubspec.yaml"}}
		updater := &dependencyUpdaterStub{err: expectedErr}
		provider := newComposite(writer, updater)

		// when
		files, err := provider.FilesChanged("/repo")

		// then
		require.ErrorIs(t, err, expectedErr)
		assert.Nil(t, files)
	})
}
