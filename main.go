package main

import (
	"os"

	"github.com/gphotosuploader/gphotos-uploader-cli/cmd"
)

func main() {
	cmd.Execute()
	os.Exit(0)
}
