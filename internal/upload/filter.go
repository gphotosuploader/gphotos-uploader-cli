package upload

import (
	"log"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/match"
)

// Filter is a Filter for file uploading
type Filter struct {
	isIncluded func(string) bool
	isExcluded func(string) bool
}

var patternDictionary = map[string][]string{
	// _ALL_FILES match with all file extensions
	"_ALL_FILES_": {"*"},

	// _ALL_VIDEO_FILES match with all video file extensions supported by Google Photos
	// Source: https://support.google.com/photos/answer/6193313.
	"_ALL_VIDEO_FILES_": {"*.mpg", "*.mod", "*.mmv", "*.tod", "*.wmv", "*.asf", "*.avi", "*.divx", "*.mov", "*.m4v", "*.3gp", "*.3g2", "*.mp4", "*.m2t", "*.m2ts", "*.mts", "*.mkv",},
}

// NewFilter returns an initialized Filter struct
func NewFilter(includePatterns []string, excludePatterns []string, allowVideos bool) *Filter {
	var f Filter

	// remove empty patterns
	includePatterns = translatePatterns(includePatterns)
	excludePatterns = translatePatterns(excludePatterns)

	if len(includePatterns) < 1 {
		includePatterns = []string{"*"}
	}

	if allowVideos {
		includePatterns = append(includePatterns, patternDictionary["_ALL_VIDEO_FILES_"]...)
	} else {
		excludePatterns = append(excludePatterns, patternDictionary["_ALL_VIDEO_FILES_"]...)
	}

	f.isIncluded = func(item string) bool {
		matched, err := match.MatchOne(includePatterns, item)
		if err != nil {
			log.Printf("error for include pattern: %v", err)
		}

		return matched
	}
	f.isExcluded = func(item string) bool {
		matched, err := match.MatchOne(excludePatterns, item)
		if err != nil {
			log.Printf("error for exclude pattern: %v", err)
		}

		return matched
	}

	return &f
}

// IsAllowed returns if an item should be uploaded.
// That means:
//   - item is in the include pattern
//   - item is not in the exclude pattern
func (f *Filter) IsAllowed(fp string) bool {
	// allow all included files that are not excluded
	return f.isIncluded(fp) && !f.isExcluded(fp)
}

// IsExcluded return if an item should be excluded.
// It's useful for skipping directories that match with an exclusion.
func (f *Filter) IsExcluded(fp string) bool {
	return f.isExcluded(fp)
}

func translatePatterns(pat []string) []string {
	var r []string
	for _, p := range pat {
		if p == "" {
			continue
		}
		hasTag := false
		for tag, val := range patternDictionary {
			if p == tag {
				r = append(r, val...)
				hasTag = true
				break
			}
		}
		if !hasTag {
			r = append(r, p)
		}
	}
	return r
}
