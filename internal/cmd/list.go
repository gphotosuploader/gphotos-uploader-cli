package cmd

import (
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd/flags"
	"github.com/spf13/cobra"
)

func NewListCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "List resources",
		Run: func(cobraCmd *cobra.Command, args []string) {
			cobraCmd.PrintErrln("Error: must also specify a resource like albums")
			cobraCmd.Usage()
		},
	}

	command.AddCommand(NewCmdListAlbums(globalFlags))

	return command
}
