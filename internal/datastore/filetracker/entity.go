package filetracker

import (
	"strings"
)

// TrackedFile represents a tracked file in the repository.
type TrackedFile struct {
	value string
}

// NewTrackedFile returns a TrackedFile with the specified values
func NewTrackedFile(value string) TrackedFile {
	return TrackedFile{
		value: value,
	}
}

// Hash returns the value value stored in the repository.
func (f TrackedFile) Hash() string {
	parts := strings.SplitN(f.value, "|", 2)

	// Previously, the value in the repository was mTime + "|" + value.
	// We decided to not use mTime in favor of the hashes. To maintain the
	// backwards compatibility, if the repository has two parts, we return
	// the second one.
	if len(parts) == 2 {
		return parts[1]
	}
	return parts[0]
}
