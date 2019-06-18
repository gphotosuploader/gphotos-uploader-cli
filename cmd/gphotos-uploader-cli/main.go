package main

import (
	"fmt"
	"os"

	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/completeduploads"
	"github.com/nmrshll/gphotos-uploader-cli/upload"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	Version string = "0.0.0"
	Build   string = "0"

	configFilePath = "~/.config/gphotos-uploader-cli/config.hjson"
)

// printErrorAndExit prints an error to stderr and force the app to exit
func printErrorAndExit(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(2)
}

func printError(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", a...)
}

func startUploader(cmd *cobra.Command, args []string) {
	uploaderConfig, err := config.LoadConfigFile(configFilePath)
	if err != nil {
		printErrorAndExit("Unable to read configuration file (%s).\nPlease review your configuration or execute 'gphotos-uploader-cli init' to create a new one.", configFilePath)
	}

	// load completedUploads DB
	db, err := leveldb.OpenFile(config.GetUploadsDBPath(), nil)
	if err != nil {
		printErrorAndExit("Error opening db: %v", err)
	}
	defer db.Close()

	// start file upload worker
	fileUploadsChan, doneUploading := upload.StartFileUploadWorker()
	doneDeleting := upload.StartDeletionsWorker()

	// launch all folder upload jobs
	for _, job := range uploaderConfig.Jobs {
		folderUploadJob := upload.NewFolderUploadJob(&job, completeduploads.NewService(db), uploaderConfig.APIAppCredentials)

		if err := folderUploadJob.Upload(fileUploadsChan); err != nil {
			printError("Failed to upload folder %s: %v", job.SourceFolder, err)
		}
	}

	// after we've run all the folder upload jobs we're done adding file upload jobs
	close(fileUploadsChan)
	// wait for all the uploads to be completed
	<-doneUploading
	fmt.Println("all uploads done")
	// after the last upload is done we're done queueing files for deletion
	upload.CloseDeletionsChan()
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
				printErrorAndExit("Failed to create the init config file: %v", err)
			}
			fmt.Printf("Configuration file has been created.\nEdit it by running:\n    nano %s\n", configFilePath)
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("gphotos-uploader-cli v%s (build: %s)\n", Version, Build)
		},
	})

	if err := rootCmd.Execute(); err != nil {
		printErrorAndExit("Could not execute the command: %v", err)
	}
}
