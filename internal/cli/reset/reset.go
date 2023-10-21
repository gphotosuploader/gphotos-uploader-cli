package reset

import (
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/flags"
	"github.com/spf13/cobra"
)

func NewCommand(globalFlags *flags.GlobalFlags) *cobra.Command {
	resetCommand := &cobra.Command{
		Use:   "reset",
		Short: "Reset internal databases",
	}

	resetCommand.AddCommand(initFileTrackerCommand(globalFlags))

	return resetCommand
}
