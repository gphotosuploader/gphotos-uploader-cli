package cmd

import (
	"context"
	"fmt"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/oauth2"

	"github.com/gphotosuploader/gphotos-uploader-cli/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/cmd/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/datastore/completeduploads"
	"github.com/gphotosuploader/gphotos-uploader-cli/datastore/leveldbstore"
	"github.com/gphotosuploader/gphotos-uploader-cli/datastore/tokenstore"
	"github.com/gphotosuploader/gphotos-uploader-cli/log"
	"github.com/gphotosuploader/gphotos-uploader-cli/photos"
	"github.com/gphotosuploader/gphotos-uploader-cli/upload"
)

// PushCmd holds the required data for the push cmd
type PushCmd struct {
	*flags.GlobalFlags
}

func NewPushCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &PushCmd{GlobalFlags: globalFlags}

	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Push local files to Google Photos service",
		Long:  `Scan configured folders in the configuration and push all new object to Google Photos service.`,
		Args:  cobra.NoArgs,
		RunE:  cmd.Run,
	}

	return pushCmd
}

func (cmd *PushCmd) Run(cobraCmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig(cmd.CfgDir)
	if err != nil {
		return fmt.Errorf("could't read configuration. Please review your configuration or run 'gphotos-uploader-cli init': file=%s, err=%s", cmd.CfgDir, err)
	}
	err = cfg.Validate()
	if err != nil {
		return fmt.Errorf("invalid configuration: file=%s, err=%s", cfg.ConfigFile(), err)
	}

	// File Tracker service to track uploaded files.
	ft, err := leveldb.OpenFile(cfg.CompletedUploadsDBDir(), nil)
	if err != nil {
		return fmt.Errorf("open completed uploads tracker failed: path=%s, err=%s", cfg.CompletedUploadsDBDir(), err)
	}
	defer ft.Close()
	fileTracker := completeduploads.NewService(completeduploads.NewLevelDBRepository(ft))

	// Token Manager service to be used as secrets backend.
	kr, err := tokenstore.NewKeyringRepository(cfg.SecretsBackendType, nil, cfg.KeyringDir())
	if err != nil {
		return fmt.Errorf("open token manager failed: type=%s, err=%s", cfg.SecretsBackendType, err)
	}
	tokenManager := tokenstore.NewService(kr)

	// Upload Session Tracker to keep upload session to resume uploads.
	uploadTracker, err := leveldbstore.NewStore(cfg.ResumableUploadsDBDir())
	if err != nil {
		return fmt.Errorf("open resumable uploads tracker failed: path=%s, err=%s", cfg.ResumableUploadsDBDir(), err)
	}
	defer uploadTracker.Close()

	// Initialize the App
	app := &app.App{
		FileTracker:   fileTracker,
		TokenManager:  tokenManager,
		UploadTracker: uploadTracker,

		Logger: log.GetInstance(),
	}

	// start file upload worker
	uploadChan, doneUploading := upload.StartFileUploadWorker(app.FileTracker, app.Logger)

	deletionQueue := upload.NewDeletionQueue()
	deletionQueue.StartWorkers(app.Logger)

	ctx := context.Background()
	// get OAuth2 Configuration with our App credentials
	oauth2Config := oauth2.Config{
		ClientID:     cfg.APIAppCredentials.ClientID,
		ClientSecret: cfg.APIAppCredentials.ClientSecret,
		Endpoint:     photos.Endpoint,
		Scopes:       photos.Scopes,
	}

	// launch all folder upload jobs
	for _, item := range cfg.Jobs {
		c, err := app.NewOAuth2Client(ctx, oauth2Config, item.Account)
		if err != nil {
			return err
		}

		gPhotos, err := gphotos.NewClientWithResumableUploads(c, app.UploadTracker)
		if err != nil {
			return err
		}

		opt := upload.NewJobOptions(item.MakeAlbums.Enabled, item.MakeAlbums.Use, item.DeleteAfterUpload, item.UploadVideos, item.IncludePatterns, item.ExcludePatterns)
		job := upload.NewFolderUploadJob(gPhotos, app.FileTracker, item.SourceFolder, opt)

		// get items to be uploaded
		jobs, err := job.ScanFolder(app.Logger)
		if err != nil {
			log.Fatalf("Failed to upload folder %s: %v", item.SourceFolder, err)
		}

		// enqueue items to be uploaded
		for _, j := range jobs {
			uploadChan <- &j
		}
	}

	// after we've run all the folder upload jobs we're done adding file upload jobs
	close(uploadChan)
	// wait for all the uploads to be completed
	<-doneUploading
	log.Done("all uploads done")

	// after the last upload is done we're done queueing files for deletion
	deletionQueue.Close()
	// wait for deletions to be completed before exiting
	deletionQueue.WaitForWorkers()
	log.Done("all deletions done")

	return nil
}
