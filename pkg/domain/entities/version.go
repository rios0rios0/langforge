package entities

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

// Version represents a semantic version value object.
type Version struct {
	raw string
	sv  *semver.Version
}

// NewVersion creates a Version from a raw version string.
// It returns an error if the string is not a valid semver.
func NewVersion(raw string) (Version, error) {
	sv, err := semver.NewVersion(raw)
	if err != nil {
		return Version{}, fmt.Errorf("invalid version %q: %w", raw, err)
	}
	return Version{raw: raw, sv: sv}, nil
}

// MustNewVersion creates a Version from a raw version string and panics on error.
func MustNewVersion(raw string) Version {
	v, err := NewVersion(raw)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns the canonical string representation.
func (v Version) String() string {
	if v.sv != nil {
		return v.sv.Original()
	}
	return v.raw
}

// SemVer returns the underlying semver.Version.
func (v Version) SemVer() *semver.Version {
	return v.sv
}

// IsZero returns true if the version is the zero value.
func (v Version) IsZero() bool {
	return v.raw == "" && v.sv == nil
}
