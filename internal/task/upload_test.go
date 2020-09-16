package task

import (
	"context"
	"errors"
	"testing"

	gphotos "github.com/gphotosuploader/google-photos-api-client-go/v2"
	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/mock"
)

func TestAlbumId(t *testing.T) {
	var testData = []struct {
		name        string
		in          string
		want        string
		errExpected bool
	}{
		{name: "WithEmptyPath", in: "", want: "", errExpected: false},
		{name: "WithSuccessfulCall", in: "bar", want: "barID", errExpected: false},
		{name: "WithFailedCall", in: "makeMyTestFail", want: "", errExpected: true},
	}

	var mockedService = &mock.GPhotosClient{
		CreateAlbumFn: func(ctx context.Context, title string) (album *photoslibrary.Album, err error) {
			if title == "makeMyTestFail" {
				return &photoslibrary.Album{}, errors.New("error")
			}
			return &photoslibrary.Album{Id: title + "ID", Title: title}, nil
		},
		AddMediaToAlbumFn: func(ctx context.Context, item gphotos.UploadItem, album *photoslibrary.Album) (*photoslibrary.MediaItem, error) {
			return &photoslibrary.MediaItem{}, nil
		},
	}

	job := EnqueuedUpload{
		PhotosClient: mockedService,
		Logger:       &mock.Logger{},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			job.AlbumName = tt.in
			got, err := job.getOrCreateAlbum()
			if got.Id != tt.want {
				t.Errorf("getOrCreateAlbumByTitle test failed: expected '%s', got '%s'", tt.want, got.Id)
			}
			if tt.errExpected && err == nil {
				t.Errorf("getOrCreateAlbumByTitle test failed: expected error")
			} else if !tt.errExpected && err != nil {
				t.Errorf("getOrCreateAlbumByTitle test failed: didn't expect error: '%s'", err)
			}
		})

	}
}
