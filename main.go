package main

import (
	"fmt"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/feedback"
	"os"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli"
)

func main() {
	gphotosCmd := cli.NewCommand()
	if err := gphotosCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error ocurred: %s", err)
		os.Exit(int(feedback.ErrGeneric))
	}
	os.Exit(int(feedback.Success))
}
