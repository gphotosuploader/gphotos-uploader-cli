package upload

import (
	"path/filepath"
	"strings"
)

// albumName returns Album name based on a path.
// If configuration option `MakeAlbum.enabled` is false, it returns empty string.
func (job *UploadFolderJob) albumName(path string) string {
	if !job.CreateAlbum {
		return ""
	}

	// AlbumName name can be customized using `MakeAlbums.use` configuration option.
	return albumNameUsingTemplate(path, job.CreateAlbumBasedOn)
}

// albumNameUsingTemplate calculate the AlbumName name for a given Path based on full folder Path (folderPath)
// or folder name (folderName).
func albumNameUsingTemplate(path, template string) string {
	switch template {
	case "folderPath":
		return albumNameUsingFolderPath(path)
	case "folderName":
		return albumNameUsingFolderName(path)
	}
	return ""
}

// albumNameUsingFolderPath returns an AlbumName name using the full Path of the given folder.
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

// albumNameUsingFolderName returns an AlbumName name using the name of the given folder.
func albumNameUsingFolderName(path string) string {
	p := filepath.Dir(path)
	if p == "." {
		return ""
	}
	return filepath.Base(p)
}
