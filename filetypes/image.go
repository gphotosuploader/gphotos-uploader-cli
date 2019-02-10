package filetypes

import (
	"fmt"
	imageLib "image"

	// register decoders for jpeg and png
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"net/http"
	"os"

	"github.com/Nr90/imgsim"
	"github.com/juju/errors"
	"github.com/steakknife/hamming"
)

func imageFromPath(filePath string) (imageLib.Image, error) {
	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err := imageLib.Decode(reader)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func imageFromURL(URL string) (imageLib.Image, error) {
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("expected http status 200, got %d", res.StatusCode)
	}

	img, _, err := imageLib.Decode(res.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// ImageTypedMedia implements TypedMedia for image files
type ImageTypedMedia struct{}

// IsCorrectlyUploaded for image file
func (im *ImageTypedMedia) IsCorrectlyUploaded(uploadedFileURL, localFilePath string) (bool, error) {
	// TODO: add sameness check for videos (use file hash) and delete if same
	if !IsImage(localFilePath) {
		return false, fmt.Errorf("%s is not an image. won't delete local file", localFilePath)
	}

	// compare uploaded image and local one
	upImg, err := imageFromURL(uploadedFileURL)
	if err != nil {
		return false, errors.Annotate(err, "failed getting image from URL")
	}
	localImg, err := imageFromPath(localFilePath)
	if err != nil {
		return false, errors.Annotate(err, "failed loading local image from path")
	}

	if isSimilarImages(upImg, localImg) {
		return true, nil
	}

	return false, nil
}

// isSimilarImages checks if two images (local and uploaded) are similar visually
// the hash used here is not a proper hash: it doesn't guarantee two images with the same hash are the same images
// it's called a perceptual hash, and can give equal or similar hashes for two different, but visually close images.
func isSimilarImages(upImg, localImg imageLib.Image) bool {
	upPerceptualHash := imgsim.DifferenceHash(upImg).String()
	localPerceptualHash := imgsim.DifferenceHash(localImg).String()

	if len(upPerceptualHash) != len(localPerceptualHash) {
		return false
	}
	hammingDistance := hamming.Strings(upPerceptualHash, localPerceptualHash)

	return hammingDistance < len(upPerceptualHash)/16
}

// // IsImageCorrectlyUploaded checks that the image that was uploaded is visually similar to the local one, before deleting the local one
// func IsImageCorrectlyUploaded(uploadedFileURL, localImgPath string) (bool, error) {
// 	if !IsImage(localImgPath) {
// 		return false, fmt.Errorf("%s is not an image. won't delete local file", localImgPath)
// 	}

// 	// compare uploaded image and local one
// 	upImg, err := imageFromURL(uploadedMediaItem.BaseUrl)
// 	if err != nil {
// 		return false, errors.Annotate(err, "failed getting image from URL")
// 	}
// 	localImg, err := imageFromPath(localImgPath)
// 	if err != nil {
// 		return false, errors.Annotate(err, "failed loading local image from path")
// 	}

// 	if isSimilarImages(upImg, localImg) {
// 		return true, nil
// 	}

// 	return false, nil
// }
