package cmd

import (
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd/flags"
	"github.com/spf13/cobra"
)

func NewListCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "List albums or media items in Google Photos",
		Run: func(cobraCmd *cobra.Command, args []string) {
			cobraCmd.PrintErrln("Error: must also specify albums or media-items")
			cobraCmd.Usage()
		},
	}

	command.AddCommand(NewCmdListAlbums(globalFlags))
	command.AddCommand(NewCmdListMediaItems(globalFlags))

	return command
}
