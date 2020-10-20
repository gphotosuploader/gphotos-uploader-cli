package filter

var patternDictionary = map[string][]string{
	// _ALL_FILES match with all file extensions
	"_ALL_FILES_": {"**"},

	// _ALL_VIDEO_FILES match with all video file extensions supported by Google Photos
	// Source: https://support.google.com/photos/answer/6193313.
	"_ALL_VIDEO_FILES_": {"**/*.mpg", "**/*.mod", "**/*.mmv", "**/*.tod", "**/*.wmv", "**/*.asf", "**/*.avi", "**/*.divx", "**/*.mov", "**/*.m4v", "**/*.3gp", "**/*.3g2", "**/*.mp4", "**/*.m2t", "**/*.m2ts", "**/*.mts", "**/*.mkv",},
}

// translatePatterns returns an array of patterns once tagged patterns has been
// resolved using patternDictionary.
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
