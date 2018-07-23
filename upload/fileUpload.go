package upload

import (
	"fmt"
	"log"

	"github.com/palantir/stacktrace"

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
		doneUploading <- struct{}{}
	}()

	return doneUploading
}

func (fileUpload *FileUpload) upload() error {
	err := fileUpload.gphotosClient.Upload(fileUpload.filePath)
	if err != nil {
		return stacktrace.Propagate(err, "failed uploading image")
	}

	fmt.Println("image uploaded successfully")

	// check phash of uploaded image
	// TODO: uncomment and fix
	// if fu.DeleteAfterUpload {
	// 	go filesHandling.CheckUploadedAndDeleteLocal(upload.URLString(), fu.filePath)
	// }
	return nil
}
