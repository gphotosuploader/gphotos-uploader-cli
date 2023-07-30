package push

import (
	"context"
	"fmt"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/uploader"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/configuration"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/feedback"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/filter"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/upload"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net/http"
)

// PushCmd holds the required data for the push cmd
type PushCmd struct {
	// command flags
	DryRunMode bool
	NoProgress bool
}

func NewCommand() *cobra.Command {
	cmd := &PushCmd{}

	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Upload local folders to Google Photos",
		Long:  `Scan local folders and upload all new objects to Google Photos.`,
		Args:  cobra.NoArgs,
		Run:   cmd.Run,
	}

	pushCmd.Flags().BoolVar(&cmd.DryRunMode, "dry-run", false, "Dry run mode")
	pushCmd.Flags().BoolVar(&cmd.DryRunMode, "no-progress", false, "Hide progress bar")

	return pushCmd
}

func (cmd *PushCmd) Run(_ *cobra.Command, _ []string) {
	ctx := context.Background()
	cli, err := app.Start(ctx, configuration.Settings.GetString("directories.data"))
	if err != nil {
		feedback.FatalError(err, feedback.ErrGeneric)
	}
	defer func() {
		_ = cli.Stop()
	}()

	type Data struct {
		Path              string   `mapstructure:"path"`
		CreateAlbums      string   `mapstructure:"create_albums,omitempty"`
		DeleteAfterUpload bool     `mapstructure:"delete_after_upload,omitempty"`
		Include           []string `mapstructure:"include,omitempty"`
		Exclude           []string `mapstructure:"exclude,omitempty"`
	}

	// launch all folder upload jobs
	var folders []Data
	if err = configuration.Settings.UnmarshalKey("folders", &folders); err != nil {
		feedback.FatalError(err, feedback.ErrCoreConfig)
	}

	// TODO: validate folders and set defaults.

	photosService, err := newPhotosService(cli.Client, cli.UploadSessionTracker, cli.Logger)
	if err != nil {
		feedback.FatalError(err, feedback.ErrGeneric)
	}

	if cmd.DryRunMode {
		logrus.Info("[DRY-RUN] Running in dry run mode. No file will be uploaded.")
	}

	for _, folder := range folders {
		sourceFolder := folder.Path

		filterFiles, err := filter.Compile(folder.Include, folder.Exclude)
		if err != nil {
			feedback.FatalError(err, feedback.ErrCoreConfig)
		}

		folder := upload.UploadFolderJob{
			FileTracker: cli.FileTracker,

			SourceFolder: sourceFolder,
			CreateAlbums: folder.CreateAlbums,
			Filter:       filterFiles,
		}

		// get UploadItem{} to be uploaded to Google Photos.
		itemsToUpload, err := folder.ScanFolder(cli.Logger)
		if err != nil {
			logrus.Fatalf("Failed to process location '%s': %s", folder.SourceFolder, err)
			continue
		}

		totalItems := len(itemsToUpload)
		var uploadedItems int

		logrus.Infof("Found %d items to be uploaded processing location '%s'.", totalItems, folder.SourceFolder)

		bar := feedback.NewTaskProgressBar("Uploading files...", totalItems, true)

		itemsGroupedByAlbum := upload.GroupByAlbum(itemsToUpload)
		for albumName, files := range itemsGroupedByAlbum {
			albumId, err := getOrCreateAlbum(ctx, photosService.Albums, albumName)
			if err != nil {
				errMsg := fmt.Sprintf("Unable to create album '%s': %s", albumName, err)
				feedback.Warning(errMsg)
				continue
			}

			for _, file := range files {
				logrus.Debugf("Processing (%d/%d): %s", uploadedItems+1, totalItems, file)

				if !cmd.DryRunMode {
					// Upload the file and add it to PhotosService.
					_, err := photosService.UploadToAlbum(ctx, albumId, file.Path)

					// Check if the Google Photos daily quota has been exceeded.
					var e *gphotos.ErrDailyQuotaExceeded
					if errors.As(err, &e) {
						feedback.Fatal("returning 'quota exceeded' error", feedback.ErrNetwork)
					}

					if err != nil {
						errMsg := fmt.Sprintf("Error processing %s: %s", file, err)
						feedback.Warning(errMsg)
						continue
					}

					// Mark the file as uploaded in the FileTracker.
					if err := cli.FileTracker.MarkAsUploaded(file.Path); err != nil {
						logrus.Warnf("Tracking file as uploaded failed: file=%s, error=%v", file, err)
					}

					//if folder.DeleteAfterUpload {
					//	if err := file.Remove(); err != nil {
					//		logrus.Errorf("Deletion request failed: file=%s, err=%v", file, err)
					//	}
					//}
				}

				bar.Add(1)
				uploadedItems++
			}
		}

		bar.Finish()

		feedback.Printf("%d processed files: %d successfully, %d with errors", totalItems, uploadedItems, totalItems-uploadedItems)
	}

	feedback.Print("All folders has been completed.")
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
