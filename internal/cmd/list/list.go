package list

import (
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd/flags"
	"github.com/spf13/cobra"
)

func NewCommand(globalFlags *flags.GlobalFlags) *cobra.Command {
	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List albums or media items in Google Photos",
	}

	listCommand.AddCommand(initAlbumsCommand(globalFlags))
	listCommand.AddCommand(initMediaItemsCommand(globalFlags))

	return listCommand
}
