//go:build unit

package ruby_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rios0rios0/langforge/pkg/infrastructure/languages/ruby"
)

func TestRubyDependencyReader(t *testing.T) {
	t.Parallel()

	t.Run("should read dependencies from Gemfile", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		gemfile := "source 'https://rubygems.org'\n\ngem 'rails', '7.0.0'\ngem 'puma', '6.0.0'\ngem 'sidekiq'\n"
		require.NoError(t, os.WriteFile(filepath.Join(dir, "Gemfile"), []byte(gemfile), 0o600))
		r := &ruby.DependencyReader{}

		// when
		deps, err := r.ReadDependencies(dir)

		// then
		require.NoError(t, err)
		require.Len(t, deps, 3)
		assert.Equal(t, "rails", deps[0].Name)
		assert.Equal(t, "7.0.0", deps[0].Current)
		assert.Equal(t, "puma", deps[1].Name)
		assert.Equal(t, "6.0.0", deps[1].Current)
		assert.Equal(t, "sidekiq", deps[2].Name)
	})

	t.Run("should read dependencies from .gemspec", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		gemspec := "Gem::Specification.new do |s|\n  s.add_runtime_dependency 'rake', '~> 13.0'\n  s.add_dependency 'bundler', '~> 2.0'\nend\n"
		require.NoError(t, os.WriteFile(filepath.Join(dir, "mygem.gemspec"), []byte(gemspec), 0o600))
		r := &ruby.DependencyReader{}

		// when
		deps, err := r.ReadDependencies(dir)

		// then
		require.NoError(t, err)
		require.Len(t, deps, 2)
		assert.Equal(t, "rake", deps[0].Name)
		assert.Equal(t, "bundler", deps[1].Name)
	})

	t.Run("should return error when no Gemfile or .gemspec found", func(t *testing.T) {
		t.Parallel()

		// given
		dir := t.TempDir()
		r := &ruby.DependencyReader{}

		// when
		_, err := r.ReadDependencies(dir)

		// then
		require.Error(t, err)
	})
}
