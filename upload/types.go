package upload

import (
	"context"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	"github.com/gphotosuploader/gphotos-uploader-cli/app"
)

// gPhotosService represents a Google Photos Service.
type gPhotosService interface {
	GetOrCreateAlbumByName(name string) (*photoslibrary.Album, error)
	AddMediaItem(ctx context.Context, path string, album string) (*photoslibrary.MediaItem, error)
}

// Job represents a job to upload all photos from the specified folder
type Job struct {
	gPhotos     gPhotosService
	fileTracker app.FileTracker

	sourceFolder string
	options      JobOptions
}

// NewFolderUploadJob creates a job based on the submitted data
func NewFolderUploadJob(gPhotos gPhotosService, fileTracker app.FileTracker, sourceFolder string, opt JobOptions) *Job {
	return &Job{
		fileTracker: fileTracker,
		gPhotos:     gPhotos,

		sourceFolder: sourceFolder,
		options:      opt,
	}
}

// JobOptions represents all the options that a job can have
type JobOptions struct {
	CreateAlbum        bool
	CreateAlbumBasedOn string
	DeleteAfterUpload  bool
	Filter             *Filter
}

// Item represents an object to be uploaded to Google Photos
type Item struct {
	gPhotos gPhotosService

	path            string
	album           string
	deleteOnSuccess bool
}
