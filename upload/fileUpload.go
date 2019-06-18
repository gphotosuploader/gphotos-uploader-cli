package upload

import (
	"log"

	"github.com/juju/errors"
	gphotos "github.com/nmrshll/google-photos-api-client-go/lib-gphotos"
	"github.com/nmrshll/gphotos-uploader-cli/filetypes"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

var (
	fileUploadsChan = make(chan *FileUpload)
)

// number of concurrent uploads
const uploadConcurrency = 5

type FileUpload struct {
	*FolderUploadJob
	filePath      string
	typedMedia    filetypes.TypedMedia
	album         *photoslibrary.Album
	gphotosClient gphotos.Client
}

// read fileUploads chan for each FileUpload struct, and upload the file to gphotos
// when the fileUploadsChan is done, signal to doneUploading
func concurrentUpload(fileUploadsChan <-chan *FileUpload, doneUploading chan<- bool) {
	semaphore := make(chan bool, uploadConcurrency)
	for fileUpload := range fileUploadsChan {
		semaphore <- true
		go func(fileUpload *FileUpload) {
			defer func() { <-semaphore }()
			err := fileUpload.upload()
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

// set up channels and start concurrentUpload
// fileUploadsChan will receive FileUpload structs and upload them
// will signal doneUploading when fileUploadsChan is done
func StartFileUploadWorker() (fileUploadsChan chan *FileUpload, doneUploading chan bool) {
	doneUploading = make(chan bool)
	fileUploadsChan = make(chan *FileUpload)
	go concurrentUpload(fileUploadsChan, doneUploading)
	return fileUploadsChan, doneUploading
}

func (fileUpload *FileUpload) upload() error {
	uploadedMediaItem, err := fileUpload.gphotosClient.UploadFile(fileUpload.filePath, fileUpload.album.Id)
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

		return QueueDeletionJob(DeletionJob{
			uploadedMediaItem.BaseUrl,
			fileUpload.filePath,
			fileUpload.typedMedia,
		})
	}
	return nil
}
