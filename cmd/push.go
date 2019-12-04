package cmd

import (
	"context"
	"fmt"
	"time"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/lib-gphotos"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/gphotosuploader/gphotos-uploader-cli/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/cmd/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/config"
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
	cfg, err := config.LoadConfigAndValidate(cmd.CfgDir)
	if err != nil {
		return fmt.Errorf("please review your configuration or run 'gphotos-uploader-cli init': file=%s, err=%s", cmd.CfgDir, err)
	}

	cli, err := app.Start(cfg)
	if err != nil {
		return err
	}
	defer func() {
		_ = cli.Stop()
	}()

	// get OAuth2 Configuration with our App credentials
	oauth2Config := oauth2.Config{
		ClientID:     cfg.APIAppCredentials.ClientID,
		ClientSecret: cfg.APIAppCredentials.ClientSecret,
		Endpoint:     photos.Endpoint,
		Scopes:       photos.Scopes,
	}

	uploadQueue := worker.NewJobQueue(cmd.NumberOfWorkers, cli.Logger)
	uploadQueue.Start()
	defer uploadQueue.Stop()
	time.Sleep(1 * time.Second) // sleeps to avoid log messages colliding with output.

	// launch all folder upload jobs
	ctx := context.Background()
	var totalItems int
	for _, config := range cfg.Jobs {
		c, err := cli.NewOAuth2Client(ctx, oauth2Config, config.Account)
		if err != nil {
			return err
		}

		photosService, err := gphotos.NewClientWithResumableUploads(c, cli.UploadTracker)
		if err != nil {
			return err
		}

		folder := upload.UploadFolderJob{
			FileTracker: cli.FileTracker,

			SourceFolder:       config.SourceFolder,
			CreateAlbum:        config.MakeAlbums.Enabled,
			CreateAlbumBasedOn: config.MakeAlbums.Use,
			Filter:             upload.NewFilter(config.IncludePatterns, config.ExcludePatterns, config.UploadVideos),
		}

		// get UploadItem{} to be uploaded to Google Photos.
		itemsToUpload, err := folder.ScanFolder(cli.Logger)
		if err != nil {
			cli.Logger.Fatalf("Failed to scan folder %s: %v", config.SourceFolder, err)
		}

		// enqueue files to be uploaded. The workers will receive it via channel.
		cli.Logger.Infof("%d files pending to be uploaded in folder '%s'.", len(itemsToUpload), config.SourceFolder)
		totalItems += len(itemsToUpload)
		for _, i := range itemsToUpload {
			uploadQueue.Submit(&upload.EnqueuedJob{
				Context:       ctx,
				PhotosService: photosService,
				FileTracker:   cli.FileTracker,
				Logger:        cli.Logger,

				Path:            i.Path,
				AlbumName:       i.AlbumName,
				DeleteOnSuccess: config.DeleteAfterUpload,
			})
		}
	}

	// get responses from the enqueued jobs
	var uploadedItems int
	for i := 0; i < totalItems; i++ {
		r := <-uploadQueue.ChanJobResults()

		if r.Err != nil {
			cli.Logger.Failf("Error processing %s", r.ID)
		} else {
			uploadedItems++
			cli.Logger.Debugf("Successfully processing %s", r.ID)
		}
	}

	cli.Logger.Infof("%d processed files: %d successfully, %d with errors", totalItems, uploadedItems, totalItems-uploadedItems)
	return nil
}
