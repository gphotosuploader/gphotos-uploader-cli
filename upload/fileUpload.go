package upload

import (
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
	"github.com/juju/errors"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/completeduploads"
	"log"
)

// number of concurrent workers uploading items
const WORKERS = 5

// Item represents an object to be uploaded to Google Photos
type Item struct {
	client *gphotos.Client

	path            string
	album           string
	deleteOnSuccess bool
}

// concurrentUpload read fileUploads chan for each Item struct, and upload the file to gphotos
// when the fileUploadsChan is done, signal to doneUploading
// TODO: We should refactor this to improve concurrency
//  eg: https://gobyexample.com/worker-pools
//  eg: https://gobyexample.com/waitgroups
//  eg: https://github.schibsted.io/spt-infrastructure/yams-delivery-images/blob/master/images/image_gif.go
func concurrentUpload(fileUploadsChan <-chan *Item, doneUploading chan<- bool, completedUploads *completeduploads.Service) {
	semaphore := make(chan bool, WORKERS)
	for fileUpload := range fileUploadsChan {
		semaphore <- true
		go func(fileUpload *Item) {
			defer func() { <-semaphore }()
			err := fileUpload.upload(completedUploads)
			if err != nil {
				log.Fatal(errors.Annotate(err, "failed uploading image"))
			}
		}(fileUpload)
	}
	// drain the semaphore
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- true
	}
	doneUploading <- true
}

// StartFileUploadWorker set up channels and start concurrentUpload
// fileUploadsChan will receive Item structs and upload them
// will signal doneUploading when fileUploadsChan is done
func StartFileUploadWorker(trackingService *completeduploads.Service) (fileUploadsChan chan *Item, doneUploading chan bool) {
	doneUploading = make(chan bool)
	fileUploadsChan = make(chan *Item)
	go concurrentUpload(fileUploadsChan, doneUploading, trackingService)
	return fileUploadsChan, doneUploading
}

// getGooglePhotosAlbumId return the Id of an album with the specified name.
// If the album doesn't exist, return an empty string.
func getGooglePhotosAlbumId(name string, c *gphotos.Client) string {
	if name == "" {
		return ""
	}

	album, err := c.GetOrCreateAlbumByName(name)
	if err != nil {
		log.Printf("error creating album: name=%s, error=%v", name, err)
		return ""
	}
	return album.Id
}

func (f *Item) upload(completedUploads *completeduploads.Service) error {
	albumId := getGooglePhotosAlbumId(f.album, f.client)
	log.Printf("uploading file: file=%s, album=%v", f.path, albumId)

	// upload the file content to Google Photos
	// TODO: Fix issue #25 - Removal of GIF & Videos is broken: https://github.com/nmrshll/gphotos-uploader-cli/issues/25
	// media, err := f.client.UploadFile(f.path, albumId)
	_, err := f.client.UploadFile(f.path, albumId)
	if err != nil {
		return errors.Annotate(err, "failed uploading image")
	}

	// mark file as uploaded in the DB
	err = completedUploads.CacheAsAlreadyUploaded(f.path)
	if err != nil {
		log.Printf("error marking file as uploaded: file=%s, error=%v", f.path, err)
	}

	// queue uploaded image for visual check of result + deletion

	// TODO: Fix issue #25 - Removal of GIF & Videos is broken: https://github.com/nmrshll/gphotos-uploader-cli/issues/25
	// v0.4.0: Disable all files removal until we fix the issue properly
	/*
		if f.deleteOnSuccess {
			// get uploaded media URL into mediaItem
			uploadedMediaItem, err := f.client.MediaItems.Get(media.Id).Do()
			if err != nil {
				return errors.Annotate(err, "failed getting uploaded mediaItem")
			}

			return QueueDeletionJob(DeletionJob{
				mediaURL: uploadedMediaItem.BaseUrl,
				mimeType: uploadedMediaItem.MimeType,
				filePath: f.path,
			})
		}
	*/
	return nil
}
