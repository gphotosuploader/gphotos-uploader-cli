package upload

import (
	"github.com/nmrshll/gphotos-uploader-cli/filetypes"
	"github.com/nmrshll/gphotos-uploader-cli/filter"
	"github.com/nmrshll/gphotos-uploader-cli/utils/filesystem"
	"log"
)

// Filter is a filter for file uploading
type Filter struct {
	isIncluded  func(string) bool
	isExcluded  func(string) bool
	allowVideos bool
}

// NewFilter returns an initialized Filter struct
func NewFilter(includePatterns []string, excludePatterns []string, allowVideos bool) *Filter {
	var f Filter

	// empty includePatterns means "*"
	for i, p := range includePatterns {
		if p == "" {
			includePatterns[i] = "*"
		}
	}
	f.isIncluded = func(item string) bool {
		matched, err := filter.List(includePatterns, item)
		if err != nil {
			log.Printf("error for include pattern: %v", err)
		}

		return matched
	}
	f.isExcluded = func(item string) bool {
		matched, err := filter.List(excludePatterns, item)
		if err != nil {
			log.Printf("error for exclude pattern: %v", err)
		}

		return matched
	}
	f.allowVideos = allowVideos

	return &f
}

// IsAllowed returns if an item should be uploaded.
// That means:
//   - item is a file
//   - item is a not a video if allowVideos is not enabled
//   - item is in the include pattern
//   - item is not in the exclude pattern
func (f *Filter) IsAllowed(fp string) bool {
	// only files are allowed
	if !filesystem.IsFile(fp) {
		return false
	}

	// check if videos are allowed
	if !f.allowVideos && filetypes.IsVideo(fp) {
		log.Printf("config doesn't allow video uploads - skipping: file=%s", fp)
		return false
	}

	// allow all included files that are not excluded
	if f.isIncluded(fp) && !f.isExcluded(fp) {
		return true
	}

	log.Printf("config doesn't allow to upload this item - skipping: file=%s", fp)
	return false

}
