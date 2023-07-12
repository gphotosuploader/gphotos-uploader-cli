package version

import (
	"github.com/gphotosuploader/gphotos-uploader-cli/version"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "Shows version number.",
		Long:  "Shows the version number of Google Photos CLI which is installed on your system.",
		Args:  cobra.NoArgs,
		Run:   runVersionCommand,
	}

	return versionCommand
}

func runVersionCommand(cmd *cobra.Command, args []string) {
	info := version.VersionInfo

	cmd.Println(info)
}
