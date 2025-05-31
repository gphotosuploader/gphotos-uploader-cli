package list

import (
	"context"
	"fmt"
	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v3"
	"github.com/gphotosuploader/google-photos-api-client-go/v3/albums"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/feedback"
	"github.com/spf13/cobra"
	"io"
	"text/tabwriter"
)

// ListAlbumsCommandOptions contains the input to the 'list albums' command.
type ListAlbumsCommandOptions struct {
	*flags.GlobalFlags

	NoHeaders  bool
	NoProgress bool
}

func initAlbumsCommand(globalFlags *flags.GlobalFlags) *cobra.Command {
	o := &ListAlbumsCommandOptions{
		GlobalFlags: globalFlags,

		NoHeaders:  false,
		NoProgress: false,
	}

	command := &cobra.Command{
		Use:   "albums",
		Short: "List albums",
		Long:  `List all the albums in Google Photos where this CLI has access to.`,
		Args:  cobra.NoArgs,
		RunE:  o.Run,
	}

	command.Flags().BoolVar(&o.NoHeaders, "no-headers", false, "Don't print the header and footer.")
	command.Flags().BoolVar(&o.NoProgress, "no-progress", false, "Don't show the progress bar.")

	return command
}

func (o *ListAlbumsCommandOptions) Run(cobraCmd *cobra.Command, args []string) error {
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

	if o.NoProgress {
		cobraCmd.Println("Getting albums from Google Photos...")
	}

	cli.Logger.Debug("Calling albums API...")

	albumsList, nextPageToken, err := photos.Albums.PaginatedList(ctx, nil)
	if err != nil {
		return err
	}

	// The progress bar is not shown when using '--no-progress' flag or in '--debug' mode.
	showProgressBar := !o.Debug && !o.NoProgress

	bar := feedback.NewTaskProgressBar("Getting albums from Google Photos...", -1, showProgressBar)

	bar.Add(len(albumsList))

	// Iterate until all pages are got
	for nextPageToken != "" {
		var response []albums.Album

		cli.Logger.Debugf("Calling albums API for page: %s", nextPageToken)

		options := &albums.PaginatedListOptions{
			PageToken: nextPageToken,
		}
		response, nextPageToken, err = photos.Albums.PaginatedList(ctx, options)
		if err != nil {
			return err
		}

		// Append current page albums to the final albums list
		albumsList = append(albumsList, response...)

		bar.Add(len(response))
	}

	bar.Finish()

	cli.Logger.Debugf("Printing album list... (%d items)", len(albumsList))

	o.printAlbumsList(albumsList, cobraCmd.OutOrStdout())

	return nil
}

func (o *ListAlbumsCommandOptions) printAlbumsList(a []albums.Album, writer io.Writer) {
	if len(a) == 0 {
		fmt.Fprintln(writer, "No albums were found!") //nolint:errcheck
		return
	}

	o.printAsTable(a, writer)
}

func (o *ListAlbumsCommandOptions) printAsTable(a []albums.Album, writer io.Writer) {
	w := tabwriter.NewWriter(writer, 0, 0, 1, ' ', 0)

	if !o.NoHeaders {
		fmt.Fprintln(w, "TITLE\t ITEMS\t ID\t") //nolint:errcheck
	}

	for _, album := range a {
		fmt.Fprintf(w, "%s\t %d\t %s\t\n", album.Title, album.TotalMediaItems, album.ID) //nolint:errcheck
	}

	if !o.NoHeaders {
		fmt.Fprintf(w, "Total: %d albums.\n", len(a)) //nolint:errcheck
	}

	w.Flush() //nolint:errcheck
}
