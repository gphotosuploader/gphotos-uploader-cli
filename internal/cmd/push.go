package cmd

import (
	"context"
	"time"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader/resumable"
	"github.com/spf13/cobra"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/filter"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/task"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/upload"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/worker"
)

// PushCmd holds the required data for the push cmd
type PushCmd struct {
	*flags.GlobalFlags

	// command flags
	NumberOfWorkers int
	DryRunMode      bool
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

	pushCmd.Flags().IntVar(&cmd.NumberOfWorkers, "workers", 1, "Number of workers")
	pushCmd.Flags().BoolVar(&cmd.DryRunMode, "dry-run", false, "Dry run mode")

	return pushCmd
}

func (cmd *PushCmd) Run(cobraCmd *cobra.Command, args []string) error {
	ctx := context.Background()
	cli, err := app.Start(ctx, cmd.CfgDir)
	if err != nil {
		return err
	}
	defer func() {
		_ = cli.Stop()
	}()

	if cmd.DryRunMode {
		cli.Logger.Info("Running in dry run mode. No changes will be made.")
	}

	uploadQueue := worker.NewJobQueue(cmd.NumberOfWorkers, cli.Logger)
	uploadQueue.Start()
	defer uploadQueue.Stop()
	time.Sleep(1 * time.Second) // sleeps to avoid log messages colliding with output.

	u, err := resumable.NewResumableUploader(cli.Client, cli.UploadSessionTracker)
	if err != nil {
		return err
	}
	photosService, err := gphotos.NewClient(cli.Client, gphotos.WithUploader(u))
	if err != nil {
		return err
	}

	// launch all folder upload jobs
	var totalItems int
	for _, config := range cli.Config.Jobs {
		srcFolder := config.SourceFolder

		folder := upload.UploadFolderJob{
			FileTracker: cli.FileTracker,

			SourceFolder:       srcFolder,
			CreateAlbum:        config.MakeAlbums.Enabled,
			CreateAlbumBasedOn: config.MakeAlbums.Use,
			Filter:             filter.New(config.IncludePatterns, config.ExcludePatterns),
		}

		// get UploadItem{} to be uploaded to Google Photos.
		itemsToUpload, err := folder.ScanFolder(cli.Logger)
		if err != nil {
			cli.Logger.Fatalf("Failed to scan folder %s: %v", config.SourceFolder, err)
		}

		cli.Logger.Infof("%d files pending to be uploaded in folder '%s'.", len(itemsToUpload), config.SourceFolder)

		// enqueue files to be uploaded. The workers will receive it via channel.
		totalItems += len(itemsToUpload)
		for _, i := range itemsToUpload {
			if cmd.DryRunMode {
				uploadQueue.Submit(&task.NoOpJob{})
			} else {
				uploadQueue.Submit(&task.EnqueuedUpload{
					Context:     ctx,
					Albums:      photosService.Albums,
					Uploads:     photosService,
					FileTracker: cli.FileTracker,
					Logger:      cli.Logger,

					Path:            i.Path,
					AlbumName:       i.AlbumName,
					DeleteOnSuccess: config.DeleteAfterUpload,
				})
			}
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

	if cmd.DryRunMode {
		cli.Logger.Info("Running in dry run mode. No changes has been made.")
	} else {
		cli.Logger.Infof("%d processed files: %d successfully, %d with errors", totalItems, uploadedItems, totalItems-uploadedItems)
	}
	return nil
}
