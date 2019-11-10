package upload

import (
	"context"
	"errors"
	"testing"

	"github.com/gphotosuploader/googlemirror/api/photoslibrary/v1"

	"github.com/gphotosuploader/gphotos-uploader-cli/log"
)

type MockGPhotosService struct {
	AlbumID string
}

func (m *MockGPhotosService) GetOrCreateAlbumByName(name string) (*photoslibrary.Album, error) {
	if name == "makeMyTestFail" {
		return &photoslibrary.Album{}, errors.New("error")
	}
	return &photoslibrary.Album{Id: m.AlbumID}, nil
}

func (m *MockGPhotosService) AddMediaItem(ctx context.Context, path string, album string) (*photoslibrary.MediaItem, error) {
	return &photoslibrary.MediaItem{}, nil
}

func TestAlbumId(t *testing.T) {
	// set an AlbumID different than expected, in order to ensure it's returning empty
	// due to CreateAlbum == false.
	job := Job{
		gPhotos: &MockGPhotosService{AlbumID: "testAlbumID"},
		options: JobOptions{
			CreateAlbum:        false,
			CreateAlbumBasedOn: "folderPath",
		},
	}

	t.Run("CreateAlbumDisabled", func(t *testing.T) {
		job.options.CreateAlbum = false
		want := ""
		got := job.albumID("foo/bar/file.jpg", log.Discard)

		if got != want {
			t.Errorf("albumID test faild: expected '%s', got '%s'", want, got)
		}
	})

	t.Run("WithEmptyPath", func(t *testing.T) {
		job.options.CreateAlbum = true
		want := ""
		got := job.albumID("", log.Discard)

		if got != want {
			t.Errorf("albumID test faild: expected '%s', got '%s'", want, got)
		}
	})

	t.Run("WithSuccessfulCall", func(t *testing.T) {
		job.options.CreateAlbum = true
		want := "testAlbumID"
		got := job.albumID("foo/bar/file.jpg", log.Discard)

		if got != want {
			t.Errorf("albumID test faild: expected '%s', got '%s'", want, got)
		}
	})

	t.Run("WithFailedCall", func(t *testing.T) {
		job.options.CreateAlbum = true
		want := ""
		got := job.albumID("makeMyTestFail/file.jpg", log.Discard)

		if got != want {
			t.Errorf("albumID test faild: expected '%s', got '%s'", want, got)
		}
	})
}

func TestAlbumNameUsingTemplate(t *testing.T) {
	var testData = []struct {
		in       string
		template string
		out      string
	}{
		{in: "foo/bar/xyz", template: "folderPath", out: "foo_bar"},
		{in: "foo/bar/xyz", template: "folderName", out: "bar"},
		{in: "foo/bar/xyz/file", template: "folderPath", out: "foo_bar_xyz"},
		{in: "foo/bar/xyz/file", template: "folderName", out: "xyz"},
		{in: "foo/bar/xyz", template: "invalidTemplate", out: ""},
	}
	for _, tt := range testData {
		got := albumNameUsingTemplate(tt.in, tt.template)
		if got != tt.out {
			t.Errorf("albumNameUsingTemplate for '%s' failed: in: '%s', expected '%s', got '%s'", tt.template, tt.in, tt.out, got)
		}
	}
}

func TestAlbumNameUsingFolderPath(t *testing.T) {
	var testData = []struct {
		in  string
		out string
	}{
		{in: "", out: ""},
		{in: "foo", out: ""},
		{in: "foo/", out: "foo"},
		{in: "foo/bar", out: "foo"},
		{in: "foo/bar/", out: "foo_bar"},
	}
	for _, tt := range testData {
		got := albumNameUsingFolderPath(tt.in)
		if got != tt.out {
			t.Errorf("albumNameUsingFolderPath for '%s' failed: expected '%s', got '%s'", tt.in, tt.out, got)
		}
	}
}

func TestAlbumNameUsingFolderName(t *testing.T) {
	var testData = []struct {
		in  string
		out string
	}{
		{in: "", out: ""},
		{in: "foo", out: ""},
		{in: "foo/", out: "foo"},
		{in: "foo/bar", out: "foo"},
		{in: "foo/bar/", out: "bar"},
	}
	for _, tt := range testData {
		got := albumNameUsingFolderName(tt.in)
		if got != tt.out {
			t.Errorf("albumNameUsingFolderName for '%s' failed: expected '%s', got '%s'", tt.in, tt.out, got)
		}
	}
}
