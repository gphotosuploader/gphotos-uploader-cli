package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/nmrshll/gphotos-uploader-cli/fileshandling"
	"github.com/nmrshll/gphotos-uploader-cli/upload"
	"github.com/spf13/cobra"
)

var (
	Version string = "0.0.0"
	Build   string = "0"

	configFilePath = "~/.config/gphotos-uploader-cli/config.hjson"
)

func startUploader(cmd *cobra.Command, args []string) {
	cfg, err := config.Load(configFilePath)
	if err != nil {
		fmt.Printf("Unable to read configuration file (%s).\nPlease review your configuration or execute 'gphotos-uploader-cli init' to create a new one.\n", configFilePath)
		os.Exit(1)
	}

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
			err := config.InitConfigFile(configFilePath)
			if err != nil {
				fmt.Printf("could not create the init config file: %v\n", err)
				os.Exit(2)
			}
			fmt.Printf("Configuration file has created.\nEdit it by running:\n    nano %s\n", configFilePath)
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
