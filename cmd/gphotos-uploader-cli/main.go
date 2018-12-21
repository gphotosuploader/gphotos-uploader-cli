package main

import (
	"fmt"
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

// errorf prints an error to stderr and force the app to exist
func errorf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(2)
}

func startUploader(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfigFile(configFilePath)
	if err != nil {
		errorf("Unable to read configuration file (%s).\nPlease review your configuration or execute 'gphotos-uploader-cli init' to create a new one.", configFilePath)
	}

	// start file upload worker
	doneUploading := upload.StartFileUploadWorker()
	doneDeleting := fileshandling.StartDeletionsWorker()

	// launch all folder upload jobs
	for _, job := range cfg.Jobs {
		folderUploadJob := upload.FolderUploadJob{FolderUploadJob: &job}
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
				errorf("Failed to create the init config file: %v", err)
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
		errorf("Could not execute the command: %v", err)
	}
}
