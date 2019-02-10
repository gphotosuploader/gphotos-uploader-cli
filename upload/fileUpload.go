package upload

import (
	"log"

	gphotos "github.com/nmrshll/google-photos-api-client-go/lib-gphotos"
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
	}
	// check upload db for previous uploads
	err = fileUpload.completedUploads.CacheAsAlreadyUploaded(fileUpload.filePath)
	if err != nil {
		log.Printf("Error marking file as uploaded: %s", fileUpload.filePath)

		// TODO: centralized logger
		// // log potentially bad images to a file
		// f, err := os.OpenFile("bad_images.log",
		// 	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		// if err != nil {
		// 	log.Println(err)
		// }
		// defer f.Close()
		// badImages := log.New(f, "", log.LstdFlags)
		// badImages.Println(fileUpload.filePath)
	}

	// queue uploaded image for visual check of result + deletion
	if fileUpload.DeleteAfterUpload {
		// get uploaded media URL into mediaItem
		uploadedMediaItem, err := fileUpload.gphotosClient.MediaItems.Get(uploadedMediaItem.Id).Do()
		if err != nil {
			return stacktrace.Propagate(err, "failed getting uploaded mediaItem")
		}

		QueueDeletionJob(uploadedMediaItem, fileUpload.filePath)
	}
	return nil
}
