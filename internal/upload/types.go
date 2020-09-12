package upload

import (
	"context"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"
)

// gPhotosService represents a Google Photos Service.
type gPhotosService interface {
	GetOrCreateAlbumByName(name string) (*photoslibrary.Album, error)
	AddMediaItem(ctx context.Context, path string, album string) (*photoslibrary.MediaItem, error)
}

// FileTracker represents a service to track already uploaded files.
type FileTracker interface {
	CacheAsAlreadyUploaded(filePath string) error
	IsAlreadyUploaded(filePath string) (bool, error)
	RemoveAsAlreadyUploaded(filePath string) error
}

// UploadFolderJob represents a job to upload all photos from the specified folder
type UploadFolderJob struct {
	FileTracker FileTracker

	SourceFolder       string
	CreateAlbum        bool
	CreateAlbumBasedOn string
	Filter             *Filter
}

// Cache represents a temporary storage to reduce number of PhotoService calls.
type Cache interface {
	Get(key string) (interface{}, error)
	Put(key string, value interface{}) error
}