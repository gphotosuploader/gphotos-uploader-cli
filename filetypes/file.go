package filetypes

import (
	"io"
	"os"

	"github.com/nmrshll/gphotos-uploader-cli/utils/filesystem"
	"gopkg.in/h2non/filetype.v1"
	filematchers "gopkg.in/h2non/filetype.v1/matchers"
)

// IsImage asserts file at filePath is an image
func IsImage(filePath string) bool {
	buf, err := filesystem.BufferHeaderFromFile(filePath, 100)
	if err != nil {
		return false
	}

	kind, _ := filetype.Image(buf)

	return kind != filetype.Unknown && kind != filematchers.TypePsd && kind != filematchers.TypeTiff && kind != filematchers.TypeCR2
}

// IsVideo asserts file at filePath is an image
func IsVideo(filePath string) bool {
	buf, err := filesystem.BufferHeaderFromFile(filePath, 100)
	if err != nil {
		return false
	}

	return filetype.IsVideo(buf)
}

// IsGif asserts file at filePath is an image
func IsGif(filePath string) bool {
	buf, err := filesystem.BufferHeaderFromFile(filePath, 100)
	if err != nil {
		return false
	}

	return filetype.IsGif(buf)
}

// IsMedia asserts file at filePath is an image or video
func IsMedia(filePath string) bool {
	return IsImage(filePath) || IsVideo(filePath) || IsGif(filePath)
}

// fileHash is a classic hash for any filetype
// it differs when one byte of the files differ
func fileHash(filePath string) (uint32, error) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer inputFile.Close()

	hasher := xxHash32.New(0xCAFE) // type hash.Hash32
	defer hasher.Reset()

	_, err = io.Copy(hasher, inputFile)
	if err != nil {
		return 0, err
	}

	return hasher.Sum32(), nil
}
