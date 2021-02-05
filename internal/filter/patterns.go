package filter

import "github.com/bmatcuk/doublestar/v2"

var patternDictionary = map[string][]string{
	// _ALL_FILES match with all file extensions
	"_ALL_FILES_": {"**"},

	// _IMAGE_EXTENSIONS_ match with the supported photos file type extensions
	// Source: https://support.google.com/photos/answer/6193313
	"_IMAGE_EXTENSIONS_": {
		"**/*.jpg", "**/*.jpeg", "**/*.png", "**/*.webp", "**/*.gif",
		"**/*.JPG", "**/*.JPEG", "**/*.PNG", "**/*.WEBP", "**/*.GIF",
	},

	// _RAW_EXTENSIONS_ match with the RAW file type extensions
	// Source: https://support.google.com/photos/answer/6193313
	// Source: https://en.wikipedia.org/wiki/Raw_image_format#Raw_filename_extensions_and_respective_camera_manufacturers
	"_RAW_EXTENSIONS_": {
		"**/*.arw", "**/*.srf", "**/*.sr2", "**/*.crw", "**/*.cr2", "**/*.cr3", "**/*.dng", "**/*.nef", "**/*.nrw", "**/*.orf", "**/*.raf", "**/*.raw", "**/*.rw2",
		"**/*.ARW", "**/*.SRF", "**/*.SR2", "**/*.CRW", "**/*.CR2", "**/*.CR3", "**/*.DNG", "**/*.NEF", "**/*.NRW", "**/*.ORF", "**/*.RAF", "**/*.RAW", "**/*.RW2",
	},

	// _ALL_VIDEO_FILES match with all video file extensions supported by Google Photos
	// Source: https://support.google.com/photos/answer/6193313.
	"_ALL_VIDEO_FILES_": {
		"**/*.mpg", "**/*.mod", "**/*.mmv", "**/*.tod", "**/*.wmv", "**/*.asf", "**/*.avi", "**/*.divx", "**/*.mov", "**/*.m4v", "**/*.3gp", "**/*.3g2", "**/*.mp4", "**/*.m2t", "**/*.m2ts", "**/*.mts", "**/*.mkv",
		"**/*.MPG", "**/*.MOD", "**/*.MMV", "**/*.TOD", "**/*.WMV", "**/*.ASF", "**/*.AVI", "**/*.DIVX", "**/*.MOV", "**/*.M4V", "**/*.3GP", "**/*.3G2", "**/*.MP4", "**/*.M2T", "**/*.M2TS", "**/*.MTS", "**/*.MKV",
	},
}

// translatePatternList returns an array of patterns once tagged patterns has been
// resolved using patternDictionary.
func translatePatternList(patternList []string) []string {
	var r []string
	for _, p := range deleteEmpty(patternList) {
		r = append(r, translatePattern(p)...)
	}
	return r
}

// translatePattern returns an array of patterns once a tagged pattern has been
// resolved using patternDictionary
func translatePattern(pattern string) []string {
	if val, exist := patternDictionary[pattern]; exist {
		return val
	}
	return []string{pattern}
}

// deleteEmpty removes empty string from an array.
func deleteEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// validatePatternList tries to validate all the patterns in the patternList and
// returns error if any of them is invalid.
func validatePatternList(patternList []string) error {
	for _, pat := range deleteEmpty(patternList) {
		if err := validatePattern(pat); err != nil {
			return err
		}
	}
	return nil
}

// validatePattern tries to use pattern and returns error if it is not valid.
func validatePattern(pattern string) error {
	_, err := doublestar.PathMatch(pattern, "x")
	return err
}

// match returns true if str matches one of the patterns. Empty patterns are ignored.
func match(patternList []string, str string) (bool, error) {
	for _, pat := range deleteEmpty(patternList) {
		matched, err := doublestar.PathMatch(pat, str)
		if err != nil {
			return false, err
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}
