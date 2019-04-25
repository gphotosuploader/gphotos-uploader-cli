package filetypes

import (
	"io"
	"os"

	"github.com/h2non/filetype"
	filematchers "github.com/h2non/filetype/matchers"
	"github.com/juju/errors"
	"github.com/nmrshll/gphotos-uploader-cli/utils/filesystem"
	"github.com/pierrec/xxHash/xxHash32"
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

// IsVideo asserts file at filePath is a video file
func IsVideo(filePath string) bool {
	buf, err := filesystem.BufferHeaderFromFile(filePath, 100)
	if err != nil {
		return false
	}

	return filetype.IsVideo(buf)
}

// IsGif asserts file at filePath is a GIF image
func IsGif(filePath string) bool {
	buf, err := filesystem.BufferHeaderFromFile(filePath, 100)
	if err != nil {
		return false
	}

	return filetype.IsMIME(buf, "image/gif")
}

// IsMedia asserts file at filePath is an image or video or gif
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

// TypedMedia has a method to check the file was correctly uploaded
type TypedMedia interface {
	IsCorrectlyUploaded(uploadedFileURL, localFilePath string) (bool, error)
}

// NewTypedMedia detects if image, gif or video and returns the appropriate TypedMedia
func NewTypedMedia(filePath string) (TypedMedia, error) {
	// process only media
	if !IsMedia(filePath) {
		return nil, errors.Errorf("failed creating typedMedia: not a media file: %s", filePath)
	}

	// if image detected
	if IsImage(filePath) && !IsGif(filePath) && !IsVideo(filePath) {
		return &ImageTypedMedia{}, nil
	}

	// if gif detected
	if IsGif(filePath) && IsImage(filePath) && !IsVideo(filePath) {
		return &GifTypedMedia{}, nil
	}

	// if video detected
	if IsVideo(filePath) && !IsImage(filePath) && !IsGif(filePath) {
		return &VideoTypedMedia{}, nil
	}

	return nil, errors.Errorf("failed creating TypedMedia from file: %s", filePath)
}
