package cmd

import (
	"fmt"
	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/completeduploads"
	"github.com/nmrshll/gphotos-uploader-cli/datastore/tokenstore"
	"github.com/nmrshll/gphotos-uploader-cli/upload"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
	"log"
	"os"
)

const defaultCfgFile = "~/.config/gphotos-uploader-cli/config.hjson"

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gphotos-uploader-cli",
	Short: "This is an unofficial Google Photos uploader",
	Long: `This application will allow you to upload your pictures to Google Photos.

You can upload folders of pictures to several Google Photos accounts and organize them in albums.

See https://github.com/nmrshll/gphotos-uploader-cli for more information.`,
	Run: startUploader,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", defaultCfgFile, "set config file")
}

func startUploader(cmd *cobra.Command, args []string) {
	var cfg *config.Config

	cfg, err := config.LoadConfigFile(cfgFile)
	if err != nil {
		log.Fatalf("Unable to read configuration file (%s).\nPlease review your configuration or execute 'gphotos-uploader-cli init' to create a new one.", cfgFile)
	}

	// load completedUploads DB
	db, err := leveldb.OpenFile(config.GetUploadsDBPath(), nil)
	if err != nil {
		log.Fatalf("Error opening db: %v", err)
	}
	defer db.Close()
	completedUploadsRepository := completeduploads.NewLevelDBRepository(db)

	// token manager service to be used as secrets backend
	kr, err := tokenstore.NewKeyringRepository(cfg.SecretsBackendType, nil)
	if err != nil {
		log.Fatalf("Unable to use the token repository: %v", err)
	}
	tkm := tokenstore.NewService(kr)

	// start file upload worker
	fileUploadsChan, doneUploading := upload.StartFileUploadWorker(completeduploads.NewService(completedUploadsRepository))
	doneDeleting := upload.StartDeletionsWorker()

	// launch all folder upload jobs
	for _, job := range cfg.Jobs {
		folderUploadJob := upload.NewFolderUploadJob(&job, completeduploads.NewService(completedUploadsRepository), cfg.APIAppCredentials, &tkm)

		if err := folderUploadJob.ScanFolder(fileUploadsChan); err != nil {
			log.Fatalf("Failed to upload folder %s: %v", job.SourceFolder, err)
		}
	}

	// after we've run all the folder upload jobs we're done adding file upload jobs
	close(fileUploadsChan)
	// wait for all the uploads to be completed
	<-doneUploading
	log.Println("all uploads done")
	// after the last upload is done we're done queueing files for deletion
	upload.CloseDeletionsChan()
	// wait for deletions to be completed before exiting
	<-doneDeleting
	log.Println("all deletions done")
}


