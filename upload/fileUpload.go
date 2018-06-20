package upload

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/simonedegiacomi/gphotosuploader/api"
	"gitlab.com/nmrshll/gphotos-uploader-go-cookies/filesHandling"
)

var (
	fileUploadsChan   = make(chan *FileUpload)
	maxPhotosToUpload = -1
)

func init() {
	startFileUploadWorker()
}

func startFileUploadWorker() {
	go func() {
		counter := 0
		for fileUpload := range fileUploadsChan {
			fileUpload.upload()
			if maxPhotosToUpload > 0 {
				counter++
				if counter >= maxPhotosToUpload {
					log.Fatal("done")
				}
			}
		}
	}()
}

type FileUpload struct {
	*FolderUploadJob
	filePath  string
	albumName string
}

func (fu *FileUpload) upload() error {
	// Open the file to upload
	file, err := os.Open(fu.filePath)
	if err != nil {
		return err
	}

	// Create an UploadOptions object that describes the upload.
	options, err := api.NewUploadOptionsFromFile(file)
	if err != nil {
		return err
	}
	if fu.albumName != "" {
		options.AlbumId = fu.albumName
	}

	// Create an upload using the NewUpload method from the api package
	upload, err := api.NewUpload(options, *fu.FolderUploadJob.Credentials)
	if err != nil {
		return err
	}

	// Finally upload the image
	err = upload.Upload()
	if err != nil {
		return err
	}
	spew.Dump(upload)
	fmt.Println("image uploaded successfully")

	// check phash of uploaded image
	if fu.DeleteAfterUpload {
		go filesHandling.CheckUploadedAndDeleteLocal(upload.URLString(), fu.filePath)
	}
	return nil
}
