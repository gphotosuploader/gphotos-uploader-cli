package upload

import (
	"path/filepath"
	"strings"
)

// albumName returns the album name based on the configured parameter.
func (job *UploadFolderJob) albumName(path string) string {
	before, after, found := strings.Cut(job.Album, ":")
	if !found {
		return ""
	}
	if before == "name" {
		return after
	}
	if before != "auto" {
		return ""
	}

	switch after {
	case "folderPath":
		return albumNameUsingFolderPath(path)
	case "folderName":
		return albumNameUsingFolderName(path)
	default:
		panic("invalid Albums parameter")
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
