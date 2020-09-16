package task

import (
	"context"
	"fmt"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/upload"
)

type GooglePĥotosService interface {
	FindAlbum(ctx context.Context, title string) (*photoslibrary.Album, error)
	CreateAlbum(ctx context.Context, title string) (*photoslibrary.Album, error)

	AddMediaToAlbum(ctx context.Context, item gphotos.UploadItem, album *photoslibrary.Album) (*photoslibrary.MediaItem, error)
}

type EnqueuedUpload struct {
	Context      context.Context
	PhotosClient GooglePĥotosService
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
	if _, err = job.PhotosClient.AddMediaToAlbum(job.Context, item, album); err != nil {
		return err
	}
	return nil
}

// getOrCreateAlbum returns the created (or existent) album in PhotosService.
func (job *EnqueuedUpload) getOrCreateAlbum() (*photoslibrary.Album, error) {
	var nullAlbum = &photoslibrary.Album{}

	// Returns if empty to avoid a PhotosService call.
	if job.AlbumName == "" {
		return nullAlbum, nil
	}

	album, err := job.PhotosClient.CreateAlbum(job.Context, job.AlbumName)
	if err != nil {
		return nullAlbum, fmt.Errorf("album creation failed: name=%s, error=%s", job.AlbumName, err)
	}
	return album, nil
}
