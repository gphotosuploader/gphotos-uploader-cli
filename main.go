package main

import (
	"os"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd"
)

func main() {
	cmd.Execute()
	os.Exit(0)
}
