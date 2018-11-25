package upload

import (
	"log"

	gphotos "github.com/nmrshll/google-photos-api-client-go/lib-gphotos"
	"github.com/nmrshll/gphotos-uploader-cli/fileshandling"
	"github.com/palantir/stacktrace"
)

var (
	fileUploadsChan = make(chan *FileUpload)
)

type FileUpload struct {
	*FolderUploadJob
	filePath      string
	albumName     string
	gphotosClient gphotos.Client
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

func (fileUpload *FileUpload) upload() error { // TODO: upload to fileUpload.AlbumName
	var albumIDVariadic []string
	if fileUpload.albumName != "" {
		album, err := fileUpload.gphotosClient.GetOrCreateAlbumByName(fileUpload.albumName)
		if err != nil {
			return stacktrace.Propagate(err, "failed GetOrCreate-ing album by name")
		}
		albumIDVariadic = append(albumIDVariadic, album.Id)
	}

	uploadedMediaItem, err := fileUpload.gphotosClient.UploadFile(fileUpload.filePath, albumIDVariadic...)
	if err != nil {
		return stacktrace.Propagate(err, "failed uploading image")
	} else {
		// check upload db for previous uploads
		err := fileshandling.MarkUploaded(fileUpload.filePath)
		if err != nil {
			log.Printf("Error marking file as uploaded: %s", fileUpload.filePath)
		}
	}

	// queue uploaded image for visual check of result + deletion
	if fileUpload.DeleteAfterUpload {
		// get uploaded media URL into mediaItem
		uploadedMediaItem, err := fileUpload.gphotosClient.MediaItems.Get(uploadedMediaItem.Id).Do()
		if err != nil {
			return stacktrace.Propagate(err, "failed getting uploaded mediaItem")
		}

		fileshandling.QueueDeletionJob(uploadedMediaItem, fileUpload.filePath)
	}
	return nil
}
