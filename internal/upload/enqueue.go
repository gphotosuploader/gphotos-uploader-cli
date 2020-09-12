package upload

import (
	"context"
	"fmt"
	"os"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

type EnqueuedJob struct {
	Context       context.Context
	PhotosService gPhotosService
	FileTracker   FileTracker
	AlbumCache    Cache
	Logger        log.Logger

	Path            string
	AlbumName       string
	DeleteOnSuccess bool
}

func (job *EnqueuedJob) Process() error {
	// Get or create the album
	albumId, err := job.albumID()
	if err != nil {
		return err
	}

	// Upload the file and add it to PhotosService.
	_, err = job.PhotosService.AddMediaItem(job.Context, job.Path, albumId)
	if err != nil {
		return err
	}

	// Mark the file as uploaded in the FileTracker.
	err = job.FileTracker.CacheAsAlreadyUploaded(job.Path)
	if err != nil {
		job.Logger.Warnf("Tracking file as uploaded failed: file=%s, error=%v", job.Path, err)
	}

	// If was requested, remove the file after being uploaded.
	if job.DeleteOnSuccess {
		if err := os.Remove(job.Path); err != nil {
			job.Logger.Errorf("Deletion request failed: file=%s, err=%v", job.Path, err)
		}
	}
	return nil
}

func (job *EnqueuedJob) ID() string {
	return job.Path
}

// albumID returns the album ID of the created (or existent) album in PhotosService.
// It uses cache to reduce number of request to PhotoService.
func (job *EnqueuedJob) albumID() (string, error) {
	// Return if empty to avoid a PhotosService call.
	if job.AlbumName == "" {
		return "", nil
	}

	albumID, err := job.AlbumCache.Get(job.AlbumName)
	if err == nil {
		return albumID.(string), nil
	}

	album, err := job.PhotosService.GetOrCreateAlbumByName(job.AlbumName)
	if err != nil {
		return "", fmt.Errorf("album creation failed: name=%s, error=%s", job.AlbumName, err)
	}

	err = job.AlbumCache.Put(job.AlbumName, album.Id)
	return album.Id, err
}
