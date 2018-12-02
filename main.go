package main

import (
	"fmt"
	"log"

	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/nmrshll/gphotos-uploader-cli/fileshandling"
	"github.com/nmrshll/gphotos-uploader-cli/upload"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	// load all config parameters
	cfg := config.Load()

	db, err := leveldb.OpenFile(config.GetUploadDBPath(), nil)
	if err != nil {
		log.Fatalf("Error opening db: %s", err)
	}
	defer db.Close()

	// start file upload worker
	doneUploading := upload.StartFileUploadWorker(db)
	doneDeleting := fileshandling.StartDeletionsWorker(db)

	// launch all folder upload jobs
	for _, job := range cfg.Jobs {
		folderUploadJob := upload.FolderUploadJob{&job}
		folderUploadJob.Run(db)
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
