package fileshandling

import (
	"fmt"
	"io/ioutil"

	"github.com/nmrshll/gphotos-uploader-cli/utils/filesystem"
	"github.com/palantir/stacktrace"
	filetype "gopkg.in/h2non/filetype.v1"
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
	buf, err := filesystem.BufferFromFile(filePath)
	if err != nil {
		return false
	}

	return filetype.IsImage(buf)
}

// IsVideo asserts file at filePath is an image
func IsVideo(filePath string) bool {
	buf, err := filesystem.BufferFromFile(filePath)
	if err != nil {
		return false
	}

	return filetype.IsVideo(buf)
}

// IsMedia asserts file at filePath is an image or video
func IsMedia(filePath string) bool {
	return IsImage(filePath) || IsVideo(filePath)
}
