package main

import (
	"fmt"

	"gitlab.com/nmrshll/gphotos-uploader-go-api/config"
	"gitlab.com/nmrshll/gphotos-uploader-go-api/fileshandling"
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
	doneDeleting := fileshandling.StartDeletionsWorker()

	// launch all folder upload jobs
	for _, job := range cfg.Jobs {
		folderUploadJob := upload.FolderUploadJob{job}
		folderUploadJob.Run()
	}
	// after we've run all the folder upload jobs we're done adding file upload jobs
	upload.CloseFileUploadsChan()
	// wait for all the uploads to be completed
	<-doneUploading
	fmt.Println("all uploads done")
	// after the last upload is done we're done queueing files for deletion
	fileshandling.CloseDeletionsChan()
	// wait for deletions to be completed before exiting
	<-doneDeleting
	fmt.Println("all deletions done")
}
