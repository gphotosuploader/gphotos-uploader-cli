package upload

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gphotosuploader/gphotos-uploader-cli/log"
)

// albumID returns the album ID if MakeAlbums is enabled
func (job *Job) albumID(path string, logger log.Logger) string {
	if !job.options.createAlbum {
		return ""
	}

	name := albumNameUsingTemplate(path, job.options.createAlbumBasedOn)
	if name == "" {
		return ""
	}

	albumID, err := job.createAlbumInGPhotos(name)
	if err != nil {
		logger.Error(err)
		return ""
	}
	return albumID
}

// albumNameUsingTemplate calculate the Album name for a given path based on full folder path (folderPath)
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

// albumNameUsingFolderPath returns an Album name using the full path of the given folder.
func albumNameUsingFolderPath(path string) string {
	p := filepath.Dir(path)
	if p == "." {
		return ""
	}
	return strings.ReplaceAll(p, "/", "_")
}

// albumNameUsingFolderName returns an Album name using the name of the given folder.
func albumNameUsingFolderName(path string) string {
	p := filepath.Dir(path)
	if p == "." {
		return ""
	}
	return filepath.Base(p)
}

// createAlbumInGPhotos returns the ID of an album with the specified name or error if fails.
// If the album didn't exist, it's created thanks to GetOrCreateAlbumByName().
func (job *Job) createAlbumInGPhotos(name string) (string, error) {
	// get album ID from Google Photos API
	album, err := job.gPhotos.GetOrCreateAlbumByName(name)
	if err != nil {
		return "", fmt.Errorf("album creation failed: name=%s, error=%s", name, err)
	}
	return album.Id, nil
}
