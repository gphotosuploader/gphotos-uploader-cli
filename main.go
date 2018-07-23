package main

import (
	"gitlab.com/nmrshll/gphotos-uploader-go-api/config"
	"gitlab.com/nmrshll/gphotos-uploader-go-api/upload"
)

const (
	imagePath = "/home/me/photos_autres/USSIS/2014_11_WE_U6/DSC_0501.JPG"
)

func main() {
	// load all config parameters
	cfg := config.Load()

	// start file upload worker
	doneUploading := upload.StartFileUploadWorker()

	// launch all folder upload jobs
	for _, job := range cfg.Jobs {
		folderUploadJob := upload.FolderUploadJob{job}
		folderUploadJob.Run()
	}
	// we're done adding file upload jobs
	upload.CloseFileUploadsChan()

	// wait for all the uploads to be complete before exiting
	<-doneUploading
}
