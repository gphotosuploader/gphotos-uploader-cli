package cmd

import (
	"context"
	"fmt"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/media_items"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd/flags"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"io"
	"text/tabwriter"
)

// ListMediaItemsOptions contains the input to the 'list media' command.
type ListMediaItemsOptions struct {
	*flags.GlobalFlags

	NoHeaders  bool
	NoProgress bool

	AlbumID string
}

// NewListMediaItemsOptions returns a ListMediaItemsOptions with defaults.
func NewListMediaItemsOptions(globalFlags *flags.GlobalFlags) *ListMediaItemsOptions {
	return &ListMediaItemsOptions{
		GlobalFlags: globalFlags,

		NoHeaders:  false,
		NoProgress: false,

		AlbumID: "",
	}
}

func NewCmdListMediaItems(globalFlags *flags.GlobalFlags) *cobra.Command {
	o := NewListMediaItemsOptions(globalFlags)

	command := &cobra.Command{
		Use:   "media-items",
		Short: "List media items",
		Long:  `List all the media items in Google Photos where this CLI has access to.`,
		Args:  cobra.NoArgs,
		RunE:  o.Run,
	}

	command.Flags().BoolVar(&o.NoHeaders, "no-headers", false, "Don't print the header and footer.")
	command.Flags().BoolVar(&o.NoProgress, "no-progress", false, "Don't show the progress bar.")
	command.Flags().StringVar(&o.AlbumID, "album-id", "", "Filter results by album ID.")

	return command
}

func (o *ListMediaItemsOptions) Run(cobraCmd *cobra.Command, args []string) error {
	ctx := context.Background()
	cli, err := app.Start(ctx, o.CfgDir)
	if err != nil {
		return err
	}
	defer func() {
		_ = cli.Stop()
	}()

	photos, err := gphotos.NewClient(cli.Client)
	if err != nil {
		return err
	}

	if o.AlbumID != "" {
		cli.Logger.Debugf("Listing media items for album ID: %s", o.AlbumID)
	}

	if o.NoProgress {
		cobraCmd.Println("Getting media items from Google Photos...")
	}

	cli.Logger.Debug("Calling media items API...")

	options := &media_items.PaginatedListOptions{
		AlbumID: o.AlbumID,
	}

	mediaItemsList, nextPageToken, err := photos.MediaItems.PaginatedList(ctx, options)
	if err != nil {
		return err
	}

	// The progress bar is not shown when using '--no-progress' flag or in '--debug' mode.
	showProgressBar := !o.Debug && !o.NoProgress

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetWriter(cobraCmd.OutOrStdout()),
		progressbar.OptionSetVisibility(showProgressBar),
		progressbar.OptionSetDescription("Getting media items from Google Photos..."),
		progressbar.OptionShowCount(),
		progressbar.OptionClearOnFinish(),
	)

	bar.Add(len(mediaItemsList))

	// Iterate until all pages are got
	for nextPageToken != "" {
		var response []media_items.MediaItem

		cli.Logger.Debugf("Calling media items API for page: %s", nextPageToken)

		options.PageToken = nextPageToken
		response, nextPageToken, err = photos.MediaItems.PaginatedList(ctx, options)
		if err != nil {
			return err
		}

		// Append current page media items to the final media items list
		mediaItemsList = append(mediaItemsList, response...)

		bar.Add(len(response))
	}

	bar.Close()

	cli.Logger.Debugf("Printing media items list... (%d items)", len(mediaItemsList))

	o.printMediaItemsList(mediaItemsList, cobraCmd.OutOrStdout())

	return nil
}

func (o *ListMediaItemsOptions) printMediaItemsList(mi []media_items.MediaItem, writer io.Writer) {
	if o.AlbumID != "" {
		fmt.Fprintf(writer, "Listing media items for album ID: %s\n", o.AlbumID)
	}

	if len(mi) == 0 {
		fmt.Fprintln(writer, "No media items were found!")
		return
	}

	o.printAsTable(mi, writer)
}

func (o *ListMediaItemsOptions) printAsTable(mi []media_items.MediaItem, writer io.Writer) {
	w := tabwriter.NewWriter(writer, 0, 0, 1, ' ', 0)

	if !o.NoHeaders {
		fmt.Fprintln(w, "FILENAME\t MIME-TYPE\t ID\t")
	}

	for _, mediaItem := range mi {
		fmt.Fprintf(w, "%s\t %s\t %s\t\n", mediaItem.Filename, mediaItem.MimeType, mediaItem.ID)
	}

	if !o.NoHeaders {
		fmt.Fprintf(w, "Total: %d media items.\n", len(mi))
	}

	w.Flush()
}
