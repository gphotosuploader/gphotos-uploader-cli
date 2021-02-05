package cmd

import (
	"context"
	"net/http"
	"time"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/google-photos-api-client-go/v2/uploader/resumable"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/filter"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
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

	uploadQueue := worker.NewJobQueue(cmd.NumberOfWorkers, cli.Logger)
	uploadQueue.Start()
	defer uploadQueue.Stop()
	time.Sleep(1 * time.Second) // sleeps to avoid log messages colliding with output.

	photosService, err := newPhotosService(cli.Client, cli.UploadSessionTracker, cli.Logger)
	if err != nil {
		return err
	}

	// launch all folder upload jobs
	var totalItems int
	for _, config := range cli.Config.Jobs {
		srcFolder := config.SourceFolder

		filterFiles, err := filter.New(config.IncludePatterns, config.ExcludePatterns)
		if err != nil {
			return err
		}

		folder := upload.UploadFolderJob{
			FileTracker: cli.FileTracker,

			SourceFolder: srcFolder,
			CreateAlbums: config.CreateAlbums,
			Filter:       filterFiles,
		}

		// get UploadItem{} to be uploaded to Google Photos.
		itemsToUpload, err := folder.ScanFolder(cli.Logger)
		if err != nil {
			cli.Logger.Fatalf("Failed to process location '%s': %s", config.SourceFolder, err)
			continue
		}

		cli.Logger.Infof("Found %d items to be uploaded processing location '%s'.", len(itemsToUpload), config.SourceFolder)

		// If dry-run-mode, stop here.
		if cmd.DryRunMode {
			cli.Logger.Info("Running in dry run mode. No changes has been made.")
			return nil
		}

		// Get items to be uploaded by album name, this reduce a lot the calls to Google Photos API to get albums ID.
		itemsByAlbum := make(map[string][]string)
		for _, i := range itemsToUpload {
			itemsByAlbum[i.AlbumName] = append(itemsByAlbum[i.AlbumName], i.Path)
		}

		for albumName, files := range itemsByAlbum {
			albumId, err := getOrCreateAlbum(ctx, photosService.Albums, albumName)
			if err != nil {
				cli.Logger.Failf("Unable to create album '%s': %s", albumName, err)
				continue
			}
			for _, file := range files {
				// enqueue files to be uploaded. The workers will receive it via channel.
				totalItems++
				uploadQueue.Submit(&task.EnqueuedUpload{
					Context:     ctx,
					Uploads:     photosService,
					FileTracker: cli.FileTracker,
					Logger:      cli.Logger,

					Path:            file,
					AlbumID:         albumId,
					DeleteOnSuccess: config.DeleteAfterUpload,
				})
			}
		}
	}

	if totalItems == 0 {
		return nil
	}

	bar := progressbar.NewOptions(totalItems,
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("Uploading files..."),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionShowCount(),
	)

	// get responses from the enqueued jobs
	var uploadedItems int
	for i := 0; i < totalItems; i++ {
		r := <-uploadQueue.ChanJobResults()

		_ = bar.Add(1)

		if r.Err != nil {
			cli.Logger.Failf("Error processing %s", r.ID)
		} else {
			uploadedItems++
			cli.Logger.Debugf("Successfully processing %s", r.ID)
		}
	}

	_ = bar.Finish()

	cli.Logger.Donef("%d processed files: %d successfully, %d with errors", totalItems, uploadedItems, totalItems-uploadedItems)
	return nil
}

func newPhotosService(client *http.Client, sessionTracker app.UploadSessionTracker, logger log.Logger) (*gphotos.Client, error) {
	u, err := resumable.NewResumableUploader(client, sessionTracker, resumable.WithLogger(logger))
	if err != nil {
		return nil, err
	}
	return gphotos.NewClient(client, gphotos.WithUploader(u))
}

// getOrCreateAlbum returns the created (or existent) album in PhotosService.
func getOrCreateAlbum(ctx context.Context, service task.AlbumsService, title string) (string, error) {
	// Returns if empty to avoid a PhotosService call.
	if title == "" {
		return "", nil
	}

	if album, err := service.GetByTitle(ctx, title); err == nil {
		return album.ID, nil
	}

	album, err := service.Create(ctx, title)
	if err != nil {
		return "", err
	}

	return album.ID, nil
}
