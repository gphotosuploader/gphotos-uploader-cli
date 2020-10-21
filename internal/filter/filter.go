package filter

import (
	"log"

	"github.com/bmatcuk/doublestar/v2"
)

// Filter is a Filter for file uploading
type Filter struct {
	isIncluded func(string) bool
	isExcluded func(string) bool
}

// New returns an initialized Filter struct
func New(includePatterns []string, excludePatterns []string) *Filter {

	includePatterns = translatePatterns(includePatterns)
	excludePatterns = translatePatterns(excludePatterns)

	if len(includePatterns) < 1 {
		includePatterns = patternDictionary["_ALL_FILES_"]
	}

	return &Filter{
		isIncluded: func(item string) bool {
			matched, err := matchAnyPattern(includePatterns, item)
			if err != nil {
				log.Printf("error for include pattern: %v", err)
			}

			return matched
		},
		isExcluded: func(item string) bool {
			matched, err := matchAnyPattern(excludePatterns, item)
			if err != nil {
				log.Printf("error for exclude pattern: %v", err)
			}

			return matched
		},
	}
}

// IsAllowed returns if an item should be uploaded.
// That means:
//   - item is in the include pattern
//   - item is not in the exclude pattern
func (f *Filter) IsAllowed(fp string) bool {
	return f.isIncluded(fp) && !f.isExcluded(fp)
}

// IsExcluded return if an item should be excluded.
// It's useful for skipping directories that match with an exclusion.
func (f *Filter) IsExcluded(fp string) bool {
	return f.isExcluded(fp)
}

// matchAnyPattern returns true if str matches one of the patterns. Empty patterns are ignored.
func matchAnyPattern(patterns []string, str string) (matched bool, err error) {
	for _, pat := range deleteEmpty(patterns) {
		matched, err := doublestar.Match(pat, str)
		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}
