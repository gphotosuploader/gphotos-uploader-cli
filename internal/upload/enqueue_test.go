package upload

import (
	"context"
	"errors"
	"testing"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/cache"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

type mockedPhotosService struct {
	albumID string
}

func (ms *mockedPhotosService) GetOrCreateAlbumByName(name string) (*photoslibrary.Album, error) {
	if name == "makeMyTestFail" {
		return &photoslibrary.Album{}, errors.New("error")
	}
	return &photoslibrary.Album{Id: ms.albumID}, nil
}

func (ms *mockedPhotosService) AddMediaItem(ctx context.Context, path string, album string) (*photoslibrary.MediaItem, error) {
	return &photoslibrary.MediaItem{}, nil
}

type mockedCache struct {
	key   string
	value interface{}
}

func (mc *mockedCache) Get(key string) (interface{}, error) {
	if key == mc.key {
		return mc.value, nil
	}
	return nil, cache.ErrNotFound
}

func (mc *mockedCache) Put(key string, value interface{}) error {
	return nil
}

func TestAlbumId(t *testing.T) {
	var testData = []struct {
		name        string
		in          string
		want        string
		errExpected bool
	}{
		{name: "WithEmptyPath", in: "", want: "", errExpected: false},
		{name: "WithAlbumInCache", in: "foo", want: "testAlbumID", errExpected: false},
		{name: "WithSuccessfulCallToPhotoService", in: "bar", want: "testAlbumID", errExpected: false},
		{name: "WithFailedCall", in: "makeMyTestFail", want: "", errExpected: true},
	}

	job := EnqueuedJob{
		PhotosService: &mockedPhotosService{albumID: "testAlbumID"},
		AlbumCache:    &mockedCache{key: "foo", value: "testAlbumID"},
		Logger:        log.Discard,
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			job.AlbumName = tt.in
			got, err := job.albumID()
			if got != tt.want {
				t.Errorf("albumID test failed: expected '%s', got '%s'", tt.want, got)
			}
			if tt.errExpected && err == nil {
				t.Errorf("albumID test failed: expected error")
			} else if !tt.errExpected && err != nil {
				t.Errorf("albumID test failed: didn't expect error: '%s'", err)
			}
		})

	}
}
