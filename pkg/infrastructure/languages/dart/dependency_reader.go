package dart

import (
	"bufio"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rios0rios0/langforge/pkg/domain/entities"
	"github.com/rios0rios0/langforge/pkg/support/fileutil"
)

// minQuotedLength is the shortest string that can carry a pair of quotes.
const minQuotedLength = 2

// sectionRe matches a top-level pubspec key, which is written flush left.
var sectionRe = regexp.MustCompile(`^([a-z_]+):\s*(?:#.*)?$`)

// entryRe matches an indented "name: value" pair. The value is optional because
// a dependency sourced from the SDK, a git remote or a path declares its source
// in a nested block instead of a version constraint.
var entryRe = regexp.MustCompile(`^(\s+)([A-Za-z_][A-Za-z0-9_]*):[ \t]*(.*)$`)

// dependencySections are the pubspec keys whose entries are packages.
//
//nolint:gochecknoglobals // read-only lookup table used as a constant
var dependencySections = map[string]struct{}{
	"dependencies":         {},
	"dev_dependencies":     {},
	"dependency_overrides": {},
}

// DependencyReader reads dependencies from pubspec.yaml.
type DependencyReader struct{}

// ReadDependencies parses pubspec.yaml and returns the declared packages.
func (r *DependencyReader) ReadDependencies(repoPath string) ([]entities.Dependency, error) {
	content, err := fileutil.ReadFile(filepath.Join(repoPath, Pubspec))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", Pubspec, err)
	}
	return ParseDependencies(content), nil
}

// ParseDependencies extracts the packages declared by a pubspec.yaml.
//
// Entries whose source is a nested block rather than a version constraint are
// skipped: an "sdk: flutter" package ships with the SDK and a git- or
// path-sourced package is pinned by revision, so in neither case is there a pub
// version an updater could compare or raise.
func ParseDependencies(content string) []entities.Dependency {
	var scan dependencyScanner
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		scan.consume(scanner.Text())
	}
	return scan.deps
}

// dependencyScanner folds a pubspec.yaml line by line into the packages it declares.
type dependencyScanner struct {
	section     string
	entryIndent string
	deps        []entities.Dependency
}

// consume folds one line into the scanner's state.
func (s *dependencyScanner) consume(line string) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return
	}

	if match := sectionRe.FindStringSubmatch(line); match != nil {
		s.section = match[1]
		s.entryIndent = ""
		return
	}

	if !s.inDependencySection() {
		return
	}

	match := entryRe.FindStringSubmatch(line)
	if match == nil {
		// A flush-left line that is not a section header ends the block.
		if !isIndented(line) {
			s.section = ""
		}
		return
	}

	s.addEntry(match[1], match[2], stripComment(match[3]))
}

// inDependencySection reports whether the current section holds packages.
func (s *dependencyScanner) inDependencySection() bool {
	_, ok := dependencySections[s.section]
	return ok
}

// addEntry records a package, unless the entry declares a nested source block
// instead of a version constraint, or is itself a line of such a block — which
// is what the deeper indentation identifies.
func (s *dependencyScanner) addEntry(indent, name, value string) {
	if s.entryIndent == "" {
		s.entryIndent = indent
	}
	if indent != s.entryIndent || value == "" {
		return
	}
	s.deps = append(s.deps, entities.NewDependency(name, unquote(value), "", Pubspec))
}

// isIndented reports whether a line is nested under a top-level key.
func isIndented(line string) bool {
	return strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t")
}

// stripComment removes a trailing comment from a constraint value.
func stripComment(value string) string {
	if idx := strings.Index(value, "#"); idx >= 0 {
		value = value[:idx]
	}
	return strings.TrimSpace(value)
}

// unquote removes the surrounding quotes pub allows around a constraint.
func unquote(value string) string {
	if len(value) >= minQuotedLength {
		first, last := value[0], value[len(value)-1]
		if first == last && (first == '\'' || first == '"') {
			return value[1 : len(value)-1]
		}
	}
	return value
}
