package versions

import (
	"context"
	"net/http"
)

// The aliases and functions below expose the unexported release types and helpers to the
// external versions_test package. The exported fetchers call fixed public URLs through
// http.DefaultClient, so the tests drive these helpers against an httptest server instead;
// as a _test.go file this shim never reaches a production build.

// GoRelease exposes goRelease to the external versions_test package.
type GoRelease = goRelease

// NodeRelease exposes nodeRelease to the external versions_test package.
type NodeRelease = nodeRelease

// EOLRelease exposes eolRelease to the external versions_test package.
type EOLRelease = eolRelease

// IsLTSRelease exposes isLTSRelease to the external versions_test package.
func IsLTSRelease(release NodeRelease) bool {
	return isLTSRelease(release)
}

// IsActiveEOL exposes isActiveEOL to the external versions_test package.
func IsActiveEOL(eol any) bool {
	return isActiveEOL(eol)
}

// FetchJSON exposes fetchJSON to the external versions_test package.
func FetchJSON(ctx context.Context, client *http.Client, url string, target any) error {
	return fetchJSON(ctx, client, url, target)
}

// FetchEndOfLifeLatest exposes fetchEndOfLifeLatest to the external versions_test package.
func FetchEndOfLifeLatest(
	ctx context.Context, client *http.Client, apiURL, name string, filter func(r EOLRelease) bool,
) (string, error) {
	return fetchEndOfLifeLatest(ctx, client, apiURL, name, filter)
}
