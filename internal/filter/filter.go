package filter

import (
	"fmt"
)

// Filter is a file filter based on allowed and excluded patterns.
type Filter struct {
	allowedList  []string
	excludedList []string
}

// New returns an initialized Filter struct. If allowedList is empty, _IMAGE_EXTENSIONS_ tagged pattern is used instead.
// It validates the patterns in allowedList and excludedList, returning error if they are not valid.
func New(allowedList []string, excludedList []string) (*Filter, error) {
	f := Filter{
		allowedList:  translatePatterns(allowedList),
		excludedList: translatePatterns(excludedList),
	}

	if len(f.allowedList) == 0 {
		f.allowedList = patternDictionary["_IMAGE_EXTENSIONS_"]
	}

	if err := f.validate(); err != nil {
		return nil, err
	}

	return &f, nil
}

// IsAllowed returns if an item is allowed.
// That means:
//   - item is in the include pattern
//   - item is not in the exclude pattern
func (f Filter) IsAllowed(fp string) bool {
	// patterns should be validated before, so no need to check error.
	matched, _ := matchAnyPattern(f.allowedList, fp)
	return matched && !f.IsExcluded(fp)
}

// IsExcluded return if an item should be excluded.
// It's useful for skipping directories that match with an exclusion.
func (f Filter) IsExcluded(fp string) bool {
	// patterns should be validated before, so no need to check error.
	matched, _ := matchAnyPattern(f.excludedList, fp)
	return matched
}

// validate returns error if allowedList or excludedList are not valid.
func (f Filter) validate() error {
	if err := validatePatterns(f.allowedList); err != nil {
		return fmt.Errorf("include patterns are invalid: %w", err)
	}
	if err := validatePatterns(f.excludedList); err != nil {
		return fmt.Errorf("exclude patterns are invalid: %w", err)
	}
	return nil
}
