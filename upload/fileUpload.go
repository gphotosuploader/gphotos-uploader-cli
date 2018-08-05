package upload

import (
	"log"

	"github.com/nmrshll/gphotos-uploader-cli/fileshandling"
	"github.com/nmrshll/gphotos-uploader-cli/gphotosapiclient"
	"github.com/palantir/stacktrace"
)

var (
	fileUploadsChan = make(chan *FileUpload)
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
		for fileUpload := range fileUploadsChan {
			err := fileUpload.upload()
			if err != nil {
				log.Fatal(stacktrace.Propagate(err, "failed uploading image"))
			}
		}
		doneUploading <- struct{}{}
	}()
	return doneUploading
}

func (fileUpload *FileUpload) upload() error {
	uploadedMediaItem, err := fileUpload.gphotosClient.UploadFile(fileUpload.filePath)
	if err != nil {
		return stacktrace.Propagate(err, "failed uploading image")
	}

	// queue uploaded image for visual check of result + deletion
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
