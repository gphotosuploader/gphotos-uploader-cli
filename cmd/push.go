package cmd

import (
	"context"
	"fmt"
	"time"

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
	"github.com/gphotosuploader/gphotos-uploader-cli/worker"
)

// PushCmd holds the required data for the push cmd
type PushCmd struct {
	*flags.GlobalFlags

	// command flags
	NumberOfWorkers int
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

	pushCmd.Flags().IntVar(&cmd.NumberOfWorkers, "workers", 5, "Number of workers")

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

	// get OAuth2 Configuration with our App credentials
	oauth2Config := oauth2.Config{
		ClientID:     cfg.APIAppCredentials.ClientID,
		ClientSecret: cfg.APIAppCredentials.ClientSecret,
		Endpoint:     photos.Endpoint,
		Scopes:       photos.Scopes,
	}

	uploadQueue := worker.NewJobQueue(cmd.NumberOfWorkers, app.Logger)
	uploadQueue.Start()
	defer uploadQueue.Stop()
	time.Sleep(1 * time.Second) // sleeps to avoid log messages colliding with output.

	// launch all folder upload jobs
	ctx := context.Background()
	for _, config := range cfg.Jobs {
		c, err := app.NewOAuth2Client(ctx, oauth2Config, config.Account)
		if err != nil {
			return err
		}

		photosService, err := gphotos.NewClientWithResumableUploads(c, app.UploadTracker)
		if err != nil {
			return err
		}

		folder := upload.UploadFolderJob{
			FileTracker: app.FileTracker,

			SourceFolder:       config.SourceFolder,
			CreateAlbum:        config.MakeAlbums.Enabled,
			CreateAlbumBasedOn: config.MakeAlbums.Use,
			Filter:             upload.NewFilter(config.IncludePatterns, config.ExcludePatterns, config.UploadVideos),
		}

		// get UploadItem{} to be uploaded to Google Photos.
		itemsToUpload, err := folder.ScanFolder(app.Logger)
		if err != nil {
			log.Fatalf("Failed to scan folder %s: %v", config.SourceFolder, err)
		}

		// enqueue files to be uploaded. The workers will receive it via channel.
		log.Infof("%d files pending to be uploaded in folder '%s'.", len(itemsToUpload), config.SourceFolder)
		for _, i := range itemsToUpload {
			uploadQueue.Submit(&upload.EnqueuedJob{
				Context:       ctx,
				PhotosService: photosService,
				FileTracker:   app.FileTracker,
				Logger:        app.Logger,

				Path:            i.Path,
				AlbumName:       i.AlbumName,
				DeleteOnSuccess: config.DeleteAfterUpload,
			})
		}
	}
	return nil
}
