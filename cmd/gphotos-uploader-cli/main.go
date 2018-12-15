package main

import (
	"fmt"
	"log"

	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/completeduploads"
	"github.com/nmrshll/gphotos-uploader-cli/fileshandling"
	"github.com/nmrshll/gphotos-uploader-cli/upload"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
)

func startUploader(cmd *cobra.Command, args []string) {
	// load all config parameters
	cfg := config.Load()

	// load completedUploads DB
	db, err := leveldb.OpenFile(config.GetUploadsDBPath(), nil)
	if err != nil {
		log.Fatalf("Error opening db: %s", err)
	}
	defer db.Close()

	// start file upload worker
	doneUploading := upload.StartFileUploadWorker()
	doneDeleting := fileshandling.StartDeletionsWorker()

	// launch all folder upload jobs
	for _, job := range cfg.Jobs {
		folderUploadJob := upload.NewFolderUploadJob(&job, completeduploads.NewService(db))
		// upload.FolderUploadJob{
		// 	&job,
		// 	completeduploads: completeduploads.NewService(db),
		// }
		// folderUploadJob.Run()
		folderUploadJob.Upload()
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

func main() {
	rootCmd := &cobra.Command{
		Use: "gphotos-uploader-cli",
		Run: startUploader,
	}
	rootCmd.AddCommand(&cobra.Command{
		Use: "init",
		Run: func(cmd *cobra.Command, args []string) {
			config.InitConfigFile()
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gphotos-uploader-cli v0.1.1")
		},
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
