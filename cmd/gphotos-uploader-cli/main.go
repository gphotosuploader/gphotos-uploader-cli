package main

import (
	"fmt"
	"log"

	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/nmrshll/gphotos-uploader-cli/fileshandling"
	"github.com/nmrshll/gphotos-uploader-cli/upload"
	"github.com/spf13/cobra"
)

var (
	BuildVersion string = "0.1.2"
	BuildRev     string = ""
)

func startUploader(cmd *cobra.Command, args []string) {
	// load all config parameters
	cfg := config.Load()

	// start file upload worker
	doneUploading := upload.StartFileUploadWorker()
	doneDeleting := fileshandling.StartDeletionsWorker()

	// launch all folder upload jobs
	for _, job := range cfg.Jobs {
		folderUploadJob := upload.FolderUploadJob{&job}
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
			fmt.Printf("gphotos-uploader-cli v%s (build: %s)\n", BuildVersion, BuildRev)
		},
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
