package main

import (
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/configuration"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/feedback"
	"os"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli"
)

func main() {
	configuration.Settings = configuration.Init(configuration.FindConfigFileInArgs(os.Args))
	gphotosCmd := cli.NewCommand()
	if err := gphotosCmd.Execute(); err != nil {
		feedback.FatalError(err, feedback.ErrGeneric)
	}
}
