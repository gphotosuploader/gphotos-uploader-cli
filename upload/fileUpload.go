package upload

import (
	"fmt"
	"log"

	"github.com/juju/errors"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
	"github.com/gphotosuploader/gphotos-uploader-cli/datastore/completeduploads"
	"github.com/gphotosuploader/gphotos-uploader-cli/datastore/uploadurls"
)

// number of concurrent workers uploading items
const maxNumberOfWorkers = 5

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
func concurrentUpload(fileUploadsChan <-chan *Item, doneUploading chan<- bool, completedUploads *completeduploads.Service, uploadURLsService *uploadurls.Service) {
	semaphore := make(chan bool, maxNumberOfWorkers)
	for fileUpload := range fileUploadsChan {
		semaphore <- true
		go func(fileUpload *Item) {
			defer func() { <-semaphore }()
			err := fileUpload.upload(completedUploads, uploadURLsService)
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
func StartFileUploadWorker(trackingService *completeduploads.Service, uploadURLsService *uploadurls.Service) (fileUploadsChan chan *Item, doneUploading chan bool) {
	doneUploading = make(chan bool)
	fileUploadsChan = make(chan *Item)
	go concurrentUpload(fileUploadsChan, doneUploading, trackingService, uploadURLsService)
	return fileUploadsChan, doneUploading
}

// getGooglePhotosAlbumID return the Id of an album with the specified name.
// If the album doesn't exist, return an empty string.
func getGooglePhotosAlbumID(name string, c *gphotos.Client) string {
	if name == "" {
		return ""
	}

	album, err := c.GetOrCreateAlbumByName(name)
	if err != nil {
		log.Printf("Album creation failed: name=%s, error=%v", name, err)
		return ""
	}
	return album.Id
}

func (f *Item) upload(completedUploads *completeduploads.Service, uploadURLsService *uploadurls.Service) error {
	albumID := getGooglePhotosAlbumID(f.album, f.client)
	log.Printf("Uploading object: file=%s", f.path)

	// check upload URL db for previous uploads
	log.Println("Looking up upload URLs database for ", f.path)
	curUploadURL, err := uploadURLsService.GetUploadURL(f.path)
	if err != nil {
		// Not found, not an error, just an empty upload URL
		curUploadURL = ""
		log.Println(err)
	}

	// upload the file content to Google Photos
	ptrUploadURL := &curUploadURL
	_, err = f.client.UploadFileResumable(f.path, ptrUploadURL, albumID)
	if err != nil {
		err = errors.Annotate(err, "failed uploading image")
	}

	if err != nil && *ptrUploadURL != "" {
		log.Printf("Error uploading file '%s', storing upload URL '%s'\n", f.path, *ptrUploadURL)
		if uploadURLsService.PutUploadURL(f.path, *ptrUploadURL) != nil {
			return fmt.Errorf("failed to store upload URL in database: %s", err)
		}

		return err
	}

	err = uploadURLsService.RemoveUploadURL(f.path)
	if err != nil {
		return fmt.Errorf("failed to remove upload URL from database: %s", err)
	}

	// mark file as uploaded in the DB
	err = completedUploads.CacheAsAlreadyUploaded(f.path)
	if err != nil {
		log.Printf("Tracking file as uploaded failed: file=%s, error=%v", f.path, err)
	}

	// queue uploaded image for visual check of result + deletion
	if f.deleteOnSuccess {
		job := DeletionJob{
			ObjectPath: f.path,
		}
		return job.Enqueue()
	}

	return nil
}
