package filetypes

import (
	"fmt"

	"github.com/juju/errors"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
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

// IsVideoCorrectlyUploaded checks that the gif that was uploaded is the same as the local one, before deleting the local one
func IsVideoCorrectlyUploaded(uploadedMediaItem *photoslibrary.MediaItem, localImgPath string) (bool, error) {
	if !IsGif(localImgPath) {
		return false, fmt.Errorf("%s is not a gif. Not deleting local file", localImgPath)
	}

	// compare uploaded image and local one
	if isSameGifs(uploadedMediaItem.BaseUrl, localImgPath) {
		return true, nil
	}

	return false, errors.New("gif was not uploaded correctly. Not deleting local file")
}
