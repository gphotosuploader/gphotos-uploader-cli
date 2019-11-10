package upload

import (
	"path/filepath"
	"strings"

	"github.com/gphotosuploader/gphotos-uploader-cli/log"
)

// albumID returns the album ID of the created (or existent) Album in Google Photos.
// If configuration option `MakeAlbum.enabled` is false, it returns empty string.
func (job *Job) albumID(path string, logger log.Logger) string {
	if !job.options.CreateAlbum {
		return ""
	}

	// Album name can be customized using `MakeAlbums.use` configuration option.
	name := albumNameUsingTemplate(path, job.options.CreateAlbumBasedOn)
	if name == "" {
		return ""
	}

	// get or creates an Album with the calculated name.
	album, err := job.gPhotos.GetOrCreateAlbumByName(name)
	if err != nil {
		logger.Errorf("album creation failed: name=%s, error=%s", name, err)
		return ""
	}
	return album.Id
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
