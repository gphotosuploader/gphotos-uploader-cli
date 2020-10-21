package filter

import (
	"fmt"
)

// Filter is a file filter based on include and exclude patterns.
type Filter struct {
	includePatterns []string
	excludePatterns []string
}

// New returns an initialized Filter struct. If includePatterns is empty, _ALL_FILES_ tagged pattern is used instead.
func New(includePatterns []string, excludePatterns []string) *Filter {
	includePatterns = translatePatterns(includePatterns)
	if len(includePatterns) == 0 {
		includePatterns = patternDictionary["_ALL_FILES_"]
	}

	f := Filter{
		includePatterns: includePatterns,
		excludePatterns: translatePatterns(excludePatterns),
	}

	return &f
}

// Validate returns error if includePatterns or excludePatterns are not valid.
func (f *Filter) Validate() error {
	if err := validatePatterns(f.includePatterns); err != nil {
		return fmt.Errorf("include patterns are invalid: %w", err)
	}
	if err := validatePatterns(f.excludePatterns); err != nil {
		return fmt.Errorf("exclude patterns are invalid: %w", err)
	}
	return nil
}

// IsAllowed returns if an item is allowed.
// That means:
//   - item is in the include pattern
//   - item is not in the exclude pattern
func (f *Filter) IsAllowed(fp string) bool {
	// patterns should be validated before, so no need to check error.
	matched, _ := matchAnyPattern(f.includePatterns, fp)
	return matched && !f.IsExcluded(fp)
}

// IsExcluded return if an item should be excluded.
// It's useful for skipping directories that match with an exclusion.
func (f *Filter) IsExcluded(fp string) bool {
	// patterns should be validated before, so no need to check error.
	matched, _ := matchAnyPattern(f.excludePatterns, fp)
	return matched
}
