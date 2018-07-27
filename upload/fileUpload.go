package upload

import (
	"fmt"
	"log"

	"github.com/palantir/stacktrace"

	"gitlab.com/nmrshll/gphotos-uploader-go-api/fileshandling"
	"gitlab.com/nmrshll/gphotos-uploader-go-api/gphotosapiclient"
)

var (
	fileUploadsChan = make(chan *FileUpload)
	// maxPhotosToUpload    = -1
)

type FileUpload struct {
	*FolderUploadJob
	filePath      string
	albumName     string
	gphotosClient gphotosapiclient.PhotosClient
}

func QueueFileUpload(fileUpload *FileUpload) {
	fileUploadsChan <- fileUpload
}
func CloseFileUploadsChan() { close(fileUploadsChan) }

func StartFileUploadWorker() (doneUploading chan struct{}) {
	doneUploading = make(chan struct{})

	go func() {
		// counter := 0
		for fileUpload := range fileUploadsChan {
			err := fileUpload.upload()
			if err != nil {
				log.Fatal(stacktrace.Propagate(err, "failed uploading image"))
			}
			// if maxPhotosToUpload > 0 {
			// 	counter++
			// 	if counter >= maxPhotosToUpload {
			// 		log.Fatal("done")
			// 	}
			// }
		}
		fmt.Println("all uploads done")
		doneUploading <- struct{}{}
	}()

	return doneUploading
}

func (fileUpload *FileUpload) upload() error {
	uploadedMediaItem, err := fileUpload.gphotosClient.UploadFile(fileUpload.filePath)
	if err != nil {
		return stacktrace.Propagate(err, "failed uploading image")
	}

	// check phash of uploaded image
	// TODO: uncomment and fix

	if fileUpload.DeleteAfterUpload {
		// get uploaded media URL into mediaItem
		uploadedMediaItem, err := fileUpload.gphotosClient.MediaItems.Get(uploadedMediaItem.Id).Do()
		if err != nil {
			return stacktrace.Propagate(err, "failed getting uploaded mediaItem")
		}

		// go fileshandling.CheckUploadedAndDeleteLocal(uploadedMediaItem, fileUpload.filePath)
		fileshandling.QueueDeletionJob(uploadedMediaItem, fileUpload.filePath)
	}
	return nil
}
