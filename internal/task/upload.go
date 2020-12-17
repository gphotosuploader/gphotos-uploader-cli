package task

import (
	"context"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/albums"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/upload"
)

type EnqueuedUpload struct {
	Context      context.Context
	PhotosClient *gphotos.Client
	FileTracker  upload.FileTracker
	Logger       log.Logger

	Path            string
	AlbumName       string
	DeleteOnSuccess bool
}

func (job *EnqueuedUpload) Process() error {
	item := upload.NewFileItem(job.Path)

	// Upload the file and add it to PhotosService.
	if err := job.addMediaToAlbum(item); err != nil {
		return err
	}

	// Mark the file as uploaded in the FileTracker.
	if err := job.FileTracker.CacheAsAlreadyUploaded(job.Path); err != nil {
		job.Logger.Warnf("Tracking file as uploaded failed: file=%s, error=%v", job.Path, err)
	}

	// If was requested, remove the file after being uploaded.
	return job.removeIfItWasRequested(item)
}

func (job *EnqueuedUpload) ID() string {
	return job.Path
}

func (job *EnqueuedUpload) removeIfItWasRequested(item upload.FileItem) error {
	if job.DeleteOnSuccess {
		if err := item.Remove(); err != nil {
			job.Logger.Errorf("Deletion request failed: file=%s, err=%v", job.Path, err)
		}
	}
	return nil
}

func (job *EnqueuedUpload) addMediaToAlbum(item upload.FileItem) error {
	// Get the album
	album, err := job.getOrCreateAlbum()
	if err != nil {
		return err
	}
	if _, err = job.PhotosClient.UploadFileToAlbum(job.Context, album.ID, item.Path); err != nil {
		return err
	}
	return nil
}

// getOrCreateAlbum returns the created (or existent) album in PhotosService.
func (job *EnqueuedUpload) getOrCreateAlbum() (*albums.Album, error) {
	// Returns if empty to avoid a PhotosService call.
	if job.AlbumName == "" {
		return &albums.Album{}, nil
	}

	if album, err := job.PhotosClient.Albums.GetByTitle(job.Context, job.AlbumName); err == nil {
		return album, nil
	}

	return  job.PhotosClient.Albums.Create(job.Context, job.AlbumName)
}
