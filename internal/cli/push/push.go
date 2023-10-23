package push

import (
	"context"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/uploader"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/feedback"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/filter"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/upload"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"net/http"
)

// PushCmd holds the required data for the push cmd
type PushCmd struct {
	*flags.GlobalFlags

	// command flags
	DryRunMode bool
}

func NewCommand(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &PushCmd{GlobalFlags: globalFlags}

	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Upload local folders to Google Photos",
		Long:  `Scan configured folders in the configuration and upload all new object to Google Photos.`,
		Args:  cobra.NoArgs,
		RunE:  cmd.Run,
	}

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

	photosService, err := newPhotosService(cli.Client, cli.UploadSessionTracker, cli.Logger)
	if err != nil {
		return err
	}

	if cmd.DryRunMode {
		cli.Logger.Info("[DRY-RUN] Running in dry run mode. No file will be uploaded.")
	}

	// launch all folder upload jobs
	for _, config := range cli.Config.Jobs {

		// TODO: CreateAlbums is maintained to ensure backwards compatibility.
		//nolint:staticcheck // I want to use deprecated method.
		if config.Album == "" && config.CreateAlbums != "" && config.CreateAlbums != "Off" {
			//nolint:staticcheck // I want to use deprecated method.
			config.Album = "auto:" + config.CreateAlbums
		}

		sourceFolder := config.SourceFolder

		filterFiles, err := filter.Compile(config.IncludePatterns, config.ExcludePatterns)
		if err != nil {
			return err
		}

		folder := upload.UploadFolderJob{
			FileTracker: cli.FileTracker,

			SourceFolder: sourceFolder,
			Album:        config.Album,
			Filter:       filterFiles,
		}

		// get UploadItem{} to be uploaded to Google Photos.
		itemsToUpload, err := folder.ScanFolder(cli.Logger)
		if err != nil {
			cli.Logger.Fatalf("Failed to process location '%s': %s", config.SourceFolder, err)
			continue
		}

		totalItems := len(itemsToUpload)
		var uploadedItems int

		cli.Logger.Infof("Found %d items to be uploaded processing location '%s'.", totalItems, config.SourceFolder)

		bar := feedback.NewTaskProgressBar("Uploading files...", totalItems, !cmd.Debug)

		itemsGroupedByAlbum := upload.GroupByAlbum(itemsToUpload)
		for albumName, files := range itemsGroupedByAlbum {
			albumId, err := getOrCreateAlbum(ctx, photosService.Albums, albumName)
			if err != nil {
				cli.Logger.Failf("Unable to create album '%s': %s", albumName, err)
				continue
			}

			for _, file := range files {
				cli.Logger.Debugf("Processing (%d/%d): %s", uploadedItems+1, totalItems, file)

				if !cmd.DryRunMode {
					// Upload the file and add it to PhotosService.
					_, err := photosService.UploadToAlbum(ctx, albumId, file.Path)

					// Check if the Google Photos daily quota has been exceeded.
					var e *gphotos.ErrDailyQuotaExceeded
					if errors.As(err, &e) {
						cli.Logger.Failf("returning 'quota exceeded' error")
						return err
					}

					if err != nil {
						cli.Logger.Failf("Error processing %s: %s", file, err)
						continue
					}

					// Mark the file as uploaded in the FileTracker.
					if err := cli.FileTracker.MarkAsUploaded(file.Path); err != nil {
						cli.Logger.Warnf("Tracking file as uploaded failed: file=%s, error=%v", file, err)
					}

					if config.DeleteAfterUpload {
						if err := file.Remove(); err != nil {
							cli.Logger.Errorf("Deletion request failed: file=%s, err=%v", file, err)
						}
					}
				}

				bar.Add(1)
				uploadedItems++
			}
		}

		bar.Finish()

		cli.Logger.Donef("%d processed files: %d successfully, %d with errors", totalItems, uploadedItems, totalItems-uploadedItems)
	}
	return nil
}

func newPhotosService(client *http.Client, sessionTracker app.UploadSessionTracker, logger log.Logger) (*gphotos.Client, error) {
	u, err := uploader.NewResumableUploader(client)
	if err != nil {
		return nil, err
	}
	u.Store = sessionTracker
	u.Logger = logger

	photos, err := gphotos.NewClient(client)
	if err != nil {
		return nil, err
	}

	// Use the resumable uploader to allow large file uploading.
	photos.Uploader = u

	return photos, nil
}

// getOrCreateAlbum returns the created (or existent) album in PhotosService.
func getOrCreateAlbum(ctx context.Context, service gphotos.AlbumsService, title string) (string, error) {
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
