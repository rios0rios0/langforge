package dart

import (
	"path/filepath"
	"regexp"

	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// Pubspec is the manifest both Dart packages and Flutter applications carry.
const Pubspec = "pubspec.yaml"

// PubspecLock is the resolved lockfile pub writes beside the manifest.
const PubspecLock = "pubspec.lock"

// flutterSectionRe matches the top-level "flutter:" key, which only a Flutter
// project declares (it configures assets, fonts and material design).
var flutterSectionRe = regexp.MustCompile(`(?m)^flutter:\s*(?:#.*)?$`)

// flutterSDKDepRe matches an "sdk: flutter" dependency source. Every Flutter
// project depends on at least one SDK-sourced package (flutter, flutter_test,
// flutter_localizations), and no pure Dart package can.
var flutterSDKDepRe = regexp.MustCompile(`(?m)^\s+sdk:\s*flutter\s*(?:#.*)?$`)

// IsFlutter reports whether the project at repoPath is a Flutter project rather
// than a plain Dart one.
//
// This is the single place that decision is made. A repository is upgraded and
// analysed with the toolchain it already uses: running `dart pub get` on a
// Flutter project fails to resolve the SDK-sourced packages, and running
// `flutter analyze` on a server-side Dart package pulls in a toolchain it does
// not need. Callers must ask this function rather than re-deriving the answer.
func IsFlutter(repoPath string) bool {
	content, err := fileutil.ReadFile(filepath.Join(repoPath, Pubspec))
	if err != nil {
		return false
	}
	return IsFlutterManifest(content)
}

// IsFlutterManifest reports whether the given pubspec.yaml content describes a
// Flutter project.
func IsFlutterManifest(content string) bool {
	return flutterSectionRe.MatchString(content) || flutterSDKDepRe.MatchString(content)
}

// executable returns the pub-capable binary for the project: "flutter" wraps
// pub for Flutter projects and "dart" provides it directly for everything else.
func executable(repoPath string) string {
	if IsFlutter(repoPath) {
		return "flutter"
	}
	return "dart"
}
