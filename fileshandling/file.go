package fileshandling

import (
	"fmt"
	"io/ioutil"

	"github.com/nmrshll/gphotos-uploader-cli/utils/filesystem"
	"github.com/palantir/stacktrace"
	filetype "gopkg.in/h2non/filetype.v1"
	filematchers "gopkg.in/h2non/filetype.v1/matchers"
	filetypes "gopkg.in/h2non/filetype.v1/types"
)

func fileBuffer(filePath string) (buf []byte, _ error) {
	if !filesystem.IsFile(filePath) {
		return nil, fmt.Errorf("not a file")
	}
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Failed finding file type: %s: Ignoring file...\n", filePath)
	}

	return buf, nil
}

// IsImage asserts file at filePath is an image
func IsImage(filePath string) bool {
	buf, err := filesystem.BufferHeaderFromFile(filePath, 100)
	if err != nil {
		return false
	}

	kind, _ := filetype.Image(buf)

	return kind != filetypes.Unknown && kind != filematchers.TypePsd && kind != filematchers.TypeTiff && kind != filematchers.TypeCR2
}

// IsVideo asserts file at filePath is an image
func IsVideo(filePath string) bool {
	buf, err := filesystem.BufferHeaderFromFile(filePath, 100)
	if err != nil {
		return false
	}

	return filetype.IsVideo(buf)
}

// IsMedia asserts file at filePath is an image or video
func IsMedia(filePath string) bool {
	return IsImage(filePath) || IsVideo(filePath)
}
