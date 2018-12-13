package fileshandling

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/palantir/stacktrace"
	filetype "gopkg.in/h2non/filetype.v1"
)

// IsFile asserts there is a file at path
func IsFile(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}

// IsDir asserts there is a directory at path
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsDir()
}

func mimeTypeContainsString(path, stringToMatch string) (bool, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return false, stacktrace.Propagate(err, "Failed finding file type: %s: Ignoring file...\n", path)
	}

	kind, err := filetype.Match(buf)
	if err != nil {
		return false, stacktrace.Propagate(err, "Failed finding file type: %s: Ignoring file...\n", path)
	}

	if strings.Contains(kind.MIME.Value, stringToMatch) {
		return true, nil
	}
	return false, nil
}

// IsImage asserts file at filePath is an image
func IsImage(filePath string) (bool, error) {
	return mimeTypeContainsString(filePath, "image")
}

// IsVideo asserts file at filePath is an image
func IsVideo(filePath string) (bool, error) {
	return mimeTypeContainsString(filePath, "video")
}

// IsMedia asserts file at filePath is an image or video
func IsMedia(filePath string) (bool, error) {
	isImage, err := IsImage(filePath)
	if isImage && err == nil {
		return isImage, err
	}

	isVideo, err := IsVideo(filePath)
	if isVideo && err == nil {
		return isVideo, err
	}

	return false, nil
}
