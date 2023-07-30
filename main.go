package main

import (
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/configuration"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/feedback"
	"os"
)

func main() {
	configuration.Settings = configuration.Init(configuration.FindConfigFileInArgs(os.Args))
	gphotosCmd := cli.New()
	if err := gphotosCmd.Execute(); err != nil {
		feedback.FatalError(err, feedback.ErrGeneric)
	}
}
