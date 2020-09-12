package upload

import (
	"context"
	"errors"
	"testing"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
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
		name         string
		in           string
		want         string
		err_expected bool
	}{
		{name: "WithEmptyPath", in: "", want: "", err_expected: false},
		{name: "WithSuccessfulCall", in: "bar", want: "testAlbumID", err_expected: false},
		{name: "WithFailedCall", in: "makeMyTestFail", want: "", err_expected: true},
	}

	job := EnqueuedUpload{
		PhotosClient: &MockedPhotosService{AlbumID: "testAlbumID"},
		Logger:       log.Discard,
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			job.AlbumName = tt.in
			got, err := job.getOrCreateAlbumByTitle()
			if got != tt.want {
				t.Errorf("getOrCreateAlbumByTitle test failed: expected '%s', got '%s'", tt.want, got)
			}
			if tt.err_expected && err == nil {
				t.Errorf("getOrCreateAlbumByTitle test failed: expected error")
			} else if !tt.err_expected && err != nil {
				t.Errorf("getOrCreateAlbumByTitle test failed: didn't expect error: '%s'", err)
			}
		})

	}
}
