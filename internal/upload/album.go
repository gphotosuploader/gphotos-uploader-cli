package upload

import (
	"path/filepath"
	"strings"
)

// albumName returns Album name based on the configured parameter.
// If configuration option is "Off" or "", it returns empty string.
func (job *UploadFolderJob) albumName(path string) string {
	switch job.CreateAlbums {
	case "Off":
		return ""
	case "folderPath":
		return albumNameUsingFolderPath(path)
	case "folderName":
		return albumNameUsingFolderName(path)
	default:
		return job.CreateAlbums
	}
}

// albumNameUsingFolderPath returns an AlbumID name using the full Path of the given folder.
func albumNameUsingFolderPath(path string) string {
	p := filepath.Dir(path)
	if p == "." {
		return ""
	}

	p = strings.ReplaceAll(p, "/", "_")

	// In path starts with '/' remove it before.
	if p[0] == '_' {
		return p[1:]
	}
	return p
}

// albumNameUsingFolderName returns an AlbumID name using the name of the given folder.
func albumNameUsingFolderName(path string) string {
	p := filepath.Dir(path)
	if p == "." {
		return ""
	}
	return filepath.Base(p)
}
