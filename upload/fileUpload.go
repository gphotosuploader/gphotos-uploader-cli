package upload

import (
	"context"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"

	"github.com/gphotosuploader/gphotos-uploader-cli/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/log"
)

// number of concurrent workers uploading items
const maxNumberOfWorkers = 5

// Item represents an object to be uploaded to Google Photos
type Item struct {
	gPhotos *gphotos.Client

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
func concurrentUpload(fileUploadsChan <-chan *Item, doneUploading chan<- bool, fileTracker app.FileTracker, log log.Logger) {
	semaphore := make(chan bool, maxNumberOfWorkers)
	for item := range fileUploadsChan {
		semaphore <- true
		go func(item *Item) {
			defer func() { <-semaphore }()
			ctx := context.TODO()
			log.Debugf("Uploading object: file=%s", item.path)
			err := item.process(ctx, fileTracker)
			if err != nil {
				log.Errorf("Failed to process media item: file=%s, err=%s", item.path, err)
				return
			}
		}(item)
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
func StartFileUploadWorker(fileTracker app.FileTracker, log log.Logger) (fileUploadsChan chan *Item, doneUploading chan bool) {
	doneUploading = make(chan bool)
	fileUploadsChan = make(chan *Item)
	go concurrentUpload(fileUploadsChan, doneUploading, fileTracker, log)
	return fileUploadsChan, doneUploading
}

func (f *Item) process(ctx context.Context, ft app.FileTracker) error {
	_, err := f.gPhotos.AddMediaItem(ctx, f.path, f.album)
	if err != nil {
		return err
	}

	err = ft.CacheAsAlreadyUploaded(f.path)
	if err != nil {
		log.Warnf("Tracking file as uploaded failed: file=%s, error=%v", f.path, err)
	}

	if f.deleteOnSuccess {
		job := DeletionJob{
			ObjectPath: f.path,
		}
		job.Enqueue()
	}
	return nil
}
