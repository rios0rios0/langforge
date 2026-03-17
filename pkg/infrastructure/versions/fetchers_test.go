//go:build unit

package versions

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestFetchLatestGoVersion(t *testing.T) {
	t.Parallel()

	t.Run("should return latest stable version when API responds with valid data", func(t *testing.T) {
		t.Parallel()

		// given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"version":"go1.23.1","stable":true},{"version":"go1.22.5","stable":true}]`))
		}))
		defer server.Close()

		// when
		var releases []goRelease
		err := fetchJSON(context.Background(), server.Client(), server.URL, &releases)
		require.NoError(t, err)

		var version string
		for _, release := range releases {
			if release.Stable {
				version = release.Version
				break
			}
		}

		// then
		assert.Equal(t, "go1.23.1", version)
	})

	t.Run("should return error when no stable version exists", func(t *testing.T) {
		t.Parallel()

		// given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"version":"go1.24rc1","stable":false}]`))
		}))
		defer server.Close()

		// when
		var releases []goRelease
		err := fetchJSON(context.Background(), server.Client(), server.URL, &releases)
		require.NoError(t, err)

		var found bool
		for _, release := range releases {
			if release.Stable {
				found = true
				break
			}
		}

		// then
		assert.False(t, found)
	})
}

func TestFetchLatestNodeVersion(t *testing.T) {
	t.Parallel()

	t.Run("should return latest LTS version when API responds with valid data", func(t *testing.T) {
		t.Parallel()

		// given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"version":"v23.0.0","lts":false},{"version":"v22.11.0","lts":"Jod"}]`))
		}))
		defer server.Close()

		// when
		var releases []nodeRelease
		err := fetchJSON(context.Background(), server.Client(), server.URL, &releases)
		require.NoError(t, err)

		var version string
		for _, release := range releases {
			if isLTSRelease(release) {
				version = release.Version
				break
			}
		}

		// then
		assert.Equal(t, "v22.11.0", version)
	})
}

func TestFetchEndOfLifeLatest(t *testing.T) {
	t.Parallel()

	t.Run("should return latest version when filter matches a release", func(t *testing.T) {
		t.Parallel()

		// given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[
				{"cycle":"3.12","latest":"3.12.4","eol":false,"lts":false},
				{"cycle":"3.11","latest":"3.11.9","eol":"2027-10-24","lts":false}
			]`))
		}))
		defer server.Close()

		// when
		version, err := fetchEndOfLifeLatest(
			context.Background(), server.Client(), server.URL, "Python",
			func(r eolRelease) bool { return isActiveEOL(r.EOL) },
		)

		// then
		require.NoError(t, err)
		assert.Equal(t, "3.12.4", version)
	})

	t.Run("should return latest LTS Java version when filter matches", func(t *testing.T) {
		t.Parallel()

		// given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[
				{"cycle":"23","latest":"23.0.1","eol":false,"lts":false},
				{"cycle":"21","latest":"21.0.4","eol":false,"lts":true},
				{"cycle":"17","latest":"17.0.12","eol":"2029-10-01","lts":true}
			]`))
		}))
		defer server.Close()

		// when
		version, err := fetchEndOfLifeLatest(
			context.Background(), server.Client(), server.URL, "Java",
			func(r eolRelease) bool { return r.LTS && isActiveEOL(r.EOL) },
		)

		// then
		require.NoError(t, err)
		assert.Equal(t, "21.0.4", version)
	})

	t.Run("should return error when no release matches the filter", func(t *testing.T) {
		t.Parallel()

		// given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`[{"cycle":"2.0","latest":"2.0.1","eol":true,"lts":false}]`))
		}))
		defer server.Close()

		// when
		_, err := fetchEndOfLifeLatest(
			context.Background(), server.Client(), server.URL, "Terraform",
			func(r eolRelease) bool { return isActiveEOL(r.EOL) },
		)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no active Terraform release found")
	})

	t.Run("should return error when server responds with non-200 status", func(t *testing.T) {
		t.Parallel()

		// given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		// when
		_, err := fetchEndOfLifeLatest(
			context.Background(), server.Client(), server.URL, "Java",
			func(r eolRelease) bool { return r.LTS && isActiveEOL(r.EOL) },
		)

		// then
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 404")
	})

	t.Run("should return error when server responds with invalid JSON", func(t *testing.T) {
		t.Parallel()

		// given
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`not json`))
		}))
		defer server.Close()

		// when
		_, err := fetchEndOfLifeLatest(
			context.Background(), server.Client(), server.URL, "Java",
			func(r eolRelease) bool { return r.LTS && isActiveEOL(r.EOL) },
		)

		// then
		require.Error(t, err)
	})
}
