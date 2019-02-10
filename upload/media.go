package upload

// import (
// 	"errors"
// 	"fmt"

// 	"github.com/nmrshll/gphotos-uploader-cli/filetypes"
// )

// type UploadableMedia interface {
// 	Upload() error
// 	IsCorrectlyUploaded() bool
// }

// // NewUploadableMedia detects if image,gif or video and returns the appropriate UploadableMedia
// func NewUploadableMedia(filePath string) (*UploadableMedia, error) {
// 	// process only media
// 	if !filetypes.isMedia(filePath) {
// 		fmt.Printf("not a media file: %s: skipping file...\n", filePath)
// 		return nil
// 	}

// 	// if image detected
// 	if filetypes.IsImage(filePath) && !filetypes.IsGif(filePath) && !filetypes.IsVideo(filePath) {
// 		return UploadableImage{}, nil
// 	}

// 	return nil, errors.New("failed creating uploadableMedia from file")
// }

// // UploadableImage implements UploadableMedia for image files
// type UploadableImage struct{}

// func (m *UploadableImage) Upload() error {

// }
