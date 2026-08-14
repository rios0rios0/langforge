package dart

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// pubspecVersionLineRe matches the top-level "version:" entry of a pubspec.yaml.
//
// The anchor is deliberate: it must sit at column 0 under (?m), because every
// other version-shaped value in a pubspec is indented — the SDK constraint under
// "environment:", the constraints under "dependencies:" and "dev_dependencies:".
// Only the package's own version is written flush left.
//
// Group 1 is the "version:" key and its trailing spaces, group 2 is the value,
// and group 3 is any trailing whitespace and comment, which callers preserve.
var pubspecVersionLineRe = regexp.MustCompile(`(?m)^(version:[ \t]*)([^\s#]+)([ \t]*(?:#.*)?)$`)

// buildSuffixRe splits a pubspec version into its semantic part and its build
// suffix. Dart writes the build number after a '+', which Flutter maps onto
// Android's versionCode and iOS's CFBundleVersion.
var buildSuffixRe = regexp.MustCompile(`^([^+]+)\+(.+)$`)

// BumpBuildNumber returns the version string to write into a pubspec.yaml whose
// version currently reads current, given the newly calculated semantic version.
//
// The semantic part always becomes newSemver. The build number is carried over
// and incremented, because Google Play and App Store Connect reject an upload
// whose build number is not greater than the previous one — dropping or freezing
// it would turn a release bump into a release that cannot ship. A manifest with
// no build number keeps none, so pub.dev packages (which discourage build
// metadata) are unaffected.
//
//	"0.1.0",         "0.2.0"  -> "0.2.0"
//	"1.0.0+1",       "1.1.0"  -> "1.1.0+2"
//	"2.10.2+021002", "2.11.0" -> "2.11.0+021003"
//	"1.0.0+nightly", "1.1.0"  -> "1.1.0+nightly"
func BumpBuildNumber(current, newSemver string) string {
	build, ok := splitBuild(current)
	if !ok {
		return newSemver
	}
	return newSemver + "+" + nextBuild(build)
}

// splitBuild returns the build suffix of a pubspec version, and whether it had one.
func splitBuild(version string) (string, bool) {
	match := buildSuffixRe.FindStringSubmatch(strings.TrimSpace(version))
	if match == nil {
		return "", false
	}
	return match[2], true
}

// nextBuild increments a numeric build suffix, preserving its zero padding.
// A non-numeric suffix (a channel name, a commit hash) is carried through
// unchanged: it is not a counter, so incrementing it would be meaningless.
func nextBuild(build string) string {
	number, err := strconv.ParseUint(build, 10, 64)
	if err != nil {
		return build
	}
	// %0*d restores the original width, so a padded counter such as "021002"
	// keeps its shape. An increment that outgrows the width simply widens.
	return fmt.Sprintf("%0*d", len(build), number+1)
}
