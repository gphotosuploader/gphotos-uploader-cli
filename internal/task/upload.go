package task

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

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

// FileUploadItem represents a local file.
type FileUploadItem string

// Open returns a stream.
// Caller should close it finally.
func (m FileUploadItem) Open() (io.ReadSeeker, int64, error) {
	f, err := os.Stat(m.String())
	if err != nil {
		return nil, 0, err
	}
	r, err := os.Open(m.String())
	if err != nil {
		return nil, 0, err
	}
	return r, f.Size(), nil
}

// Name returns the filename.
func (m FileUploadItem) Name() string {
	return path.Base(m.String())
}

func (m FileUploadItem) String() string {
	return string(m)
}

func (m FileUploadItem) Size() int64 {
	f, err := os.Stat(m.String())
	if err != nil {
		return 0
	}
	return f.Size()
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
	// Get or create the album
	album, err := job.getOrCreateAlbumByTitle()
	if err != nil {
		return err
	}

	// Upload the file and add it to PhotosService.
	_, err = job.PhotosClient.AddMediaToAlbum(job.Context, FileUploadItem(job.Path), album)
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

func (job *EnqueuedUpload) ID() string {
	return job.Path
}

// getOrCreateAlbumByTitle returns the album ID of the created (or existent) album in PhotosService.
func (job *EnqueuedUpload) getOrCreateAlbumByTitle() (*photoslibrary.Album, error) {
	// Return if empty to avoid a PhotosService call.
	if job.AlbumName == "" {
		return &photoslibrary.Album{}, nil
	}

	album, err := job.PhotosClient.CreateAlbum(job.Context, job.AlbumName)
	if err != nil {
		return &photoslibrary.Album{}, fmt.Errorf("Album creation failed: name=%s, error=%s", job.AlbumName, err)
	}
	return album, nil
}
