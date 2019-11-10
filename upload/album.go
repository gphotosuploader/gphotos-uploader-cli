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

	var name string
	switch job.options.createAlbumBasedOn {
	case "folderPath":
		name = strings.ReplaceAll(filepath.Dir(path), "/", "_")
	case "folderName":
	default:
		name = filepath.Base(filepath.Dir(path))
	}

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
