package upload

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlbumName(t *testing.T) {
	var testData = []struct {
		name         string
		createAlbums string

		in   string
		want string
	}{
		{
			name:         "createAlbumDisabled_With_Off",
			createAlbums: "Off",
			in:           "/foo/bar/file.jpg",
			want:         "",
		},
		{
			name:         "createAlbum_With_folderName",
			createAlbums: "folderName",
			in:           "/foo/bar/file.jpg",
			want:         "bar",
		},
		{
			name:         "createAlbum_With_folderPath",
			createAlbums: "folderPath",
			in:           "/foo/bar/file.jpg",
			want:         "foo_bar",
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			job := UploadFolderJob{
				CreateAlbums: tt.createAlbums,
			}

			assert.Equal(t, tt.want, job.albumName(tt.in))
		})

	}
}

func TestAlbumNameWithInvalidParameter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("A Panic was expected but not reached.")
		}
	}()
	job := UploadFolderJob{
		CreateAlbums: "FooBar",
	}
	_ = job.albumName("/foo/bar/file.jpg")
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
		{in: "/foo/bar", out: "foo"},
		{in: "/foo/bar/", out: "foo_bar"},
	}
	for _, tt := range testData {
		assert.Equal(t, tt.out, albumNameUsingFolderPath(tt.in))
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
		{in: "/foo/bar", out: "foo"},
		{in: "/foo/bar/", out: "bar"},
	}
	for _, tt := range testData {
		assert.Equal(t, tt.out, albumNameUsingFolderName(tt.in))
	}
}
