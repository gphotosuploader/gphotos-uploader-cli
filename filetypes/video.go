package filetypes

import (
	"fmt"

	"github.com/juju/errors"
)

// isSameVideos checks if two gifs (local and uploaded) are exactly the same
func isSameVideos(upGifPath, localGifPath string) bool {
	upHash, err := fileHash(upGifPath)
	if err != nil {
		return false
	}
	localHash, err := fileHash(localGifPath)
	if err != nil {
		return false
	}

	return upHash == localHash
}

// VideoTypedMedia implements TypedMedia for video files
type VideoTypedMedia struct{}

// IsCorrectlyUploaded checks that the video that was uploaded is the same as the local one, before deleting the local one
func (gm *VideoTypedMedia) IsCorrectlyUploaded(uploadedFileURL, localFilePath string) (bool, error) {
	if !IsGif(localFilePath) {
		return false, fmt.Errorf("%s is not a gif. Not deleting local file", localFilePath)
	}

	// compare uploaded image and local one
	if isSameGifs(uploadedFileURL, localFilePath) {
		return true, nil
	}

	return false, errors.New("gif was not uploaded correctly. Not deleting local file")
}
