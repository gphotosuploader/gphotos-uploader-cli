package upload

import (
	"context"
	"errors"
	"testing"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	"github.com/gphotosuploader/gphotos-uploader-cli/log"
)

type MockedPhotosService struct {
	AlbumID string
}

func (m *MockedPhotosService) GetOrCreateAlbumByName(name string) (*photoslibrary.Album, error) {
	if name == "makeMyTestFail" {
		return &photoslibrary.Album{}, errors.New("error")
	}
	return &photoslibrary.Album{Id: m.AlbumID}, nil
}

func (m *MockedPhotosService) AddMediaItem(ctx context.Context, path string, album string) (*photoslibrary.MediaItem, error) {
	return &photoslibrary.MediaItem{}, nil
}

func TestAlbumId(t *testing.T) {
	var testData = []struct {
		name string
		in   string
		want string
	}{
		{name: "WithEmptyPath", in: "", want: ""},
		{name: "WithSuccessfulCall", in: "bar", want: "testAlbumID"},
		{name: "WithFailedCall", in: "makeMyTestFail", want: ""},
	}

	job := EnqueuedJob{
		PhotosService: &MockedPhotosService{AlbumID: "testAlbumID"},
		Logger:        log.Discard,
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			job.AlbumName = tt.in
			got, _ := job.albumID()
			if got != tt.want {
				t.Errorf("albumID test failed: expected '%s', got '%s'", tt.want, got)
			}
		})

	}
}
