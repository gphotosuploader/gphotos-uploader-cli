package list

import (
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	listCommand := &cobra.Command{
		Use:   "list",
		Short: "List albums or media items in Google Photos",
	}

	listCommand.AddCommand(initAlbumsCommand())
	listCommand.AddCommand(initMediaItemsCommand())

	return listCommand
}
