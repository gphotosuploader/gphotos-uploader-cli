package upload

import (
	"log"

	"github.com/juju/errors"
	gphotos "github.com/nmrshll/google-photos-api-client-go/lib-gphotos"
	"github.com/nmrshll/gphotos-uploader-cli/filetypes"
)

var (
	fileUploadsChan = make(chan *FileUpload)
)

type FileUpload struct {
	*FolderUploadJob
	filePath      string
	typedMedia    filetypes.TypedMedia
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
				log.Fatal(errors.Annotate(err, "failed uploading image"))
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
			return errors.Annotate(err, "failed GetOrCreate-ing album by name")
		}
		albumIDVariadic = append(albumIDVariadic, album.Id)
	}

	uploadedMediaItem, err := fileUpload.gphotosClient.UploadFile(fileUpload.filePath, albumIDVariadic...)
	if err != nil {
		return errors.Annotate(err, "failed uploading image")
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
			return errors.Annotate(err, "failed getting uploaded mediaItem")
		}

		QueueDeletionJob(DeletionJob{
			uploadedMediaItem.BaseUrl,
			fileUpload.filePath,
			fileUpload.typedMedia,
		})
	}
	return nil
}
