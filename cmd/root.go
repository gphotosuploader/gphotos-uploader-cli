package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/gphotosuploader/gphotos-uploader-cli/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/datastore/completeduploads"
	"github.com/gphotosuploader/gphotos-uploader-cli/datastore/tokenstore"
	"github.com/gphotosuploader/gphotos-uploader-cli/upload"
)

const defaultCfgDir = "~/.config/gphotos-uploader-cli"

var cfgDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gphotos-uploader-cli",
	Short: "This is an unofficial Google Photos uploader",
	Long: `This application will allow you to upload your pictures to Google Photos.

You can upload folders of pictures to several Google Photos accounts and organize them in albums.

See https://github.com/gphotosuploader/gphotos-uploader-cli for more information.`,
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
	rootCmd.PersistentFlags().StringVar(&cfgDir, "config", defaultCfgDir, "set config folder")
}

func startUploader(cmd *cobra.Command, args []string) {
	var cfg *config.Config

	cfg, err := config.LoadConfigFile(cfgDir)
	if err != nil {
		log.Fatalf("Unable to read configuration from '%s'.\nPlease review your configuration or execute 'gphotos-uploader-cli init' to create a new one.", cfgDir)
	}

	// load completedUploads DB
	db, err := leveldb.OpenFile(cfg.TrackingDBPath, nil)
	if err != nil {
		log.Fatalf("Error opening db: path=%s, err=%v", cfg.TrackingDBPath, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	fileTracker := completeduploads.NewService(completeduploads.NewLevelDBRepository(db))

	// token manager service to be used as secrets backend
	kr, err := tokenstore.NewKeyringRepository(cfg.SecretsBackendType, nil)
	if err != nil {
		log.Fatalf("Unable to use the token repository: %v", err)
	}
	tkm := tokenstore.NewService(kr)

	// start file upload worker
	uploadChan, doneUploading := upload.StartFileUploadWorker(fileTracker)

	deletionQueue := upload.NewDeletionQueue()
	deletionQueue.StartWorkers()

	// get OAuth2 Configuration with our App credentials
	oauthConfig := config.OAuthConfig(cfg.APIAppCredentials)

	// launch all folder upload jobs
	for _, item := range cfg.Jobs {
		gPhotos, err := upload.Authenticate(tkm, oauthConfig, item.Account)
		if err != nil {
			log.Fatal(err)
		}

		opt := upload.NewJobOptions(item.MakeAlbums.Enabled, item.DeleteAfterUpload, item.UploadVideos, item.IncludePatterns, item.ExcludePatterns)
		job := upload.NewFolderUploadJob(&gPhotos.Client, fileTracker, item.SourceFolder, opt)

		if err := job.ScanFolder(uploadChan); err != nil {
			log.Fatalf("Failed to upload folder %s: %v", item.SourceFolder, err)
		}
	}

	// after we've run all the folder upload jobs we're done adding file upload jobs
	close(uploadChan)
	// wait for all the uploads to be completed
	<-doneUploading
	log.Println("all uploads done")

	// after the last upload is done we're done queueing files for deletion
	deletionQueue.Close()
	// wait for deletions to be completed before exiting
	deletionQueue.WaitForWorkers()
	log.Println("all deletions done")
}
