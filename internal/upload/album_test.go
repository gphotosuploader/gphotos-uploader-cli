package upload

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlbumName(t *testing.T) {
	var testData = []struct {
		name  string
		album string

		in   string
		want string
	}{
		{
			name:  "album set an album's name",
			album: "name:albumName",
			in:    "/foo/bar/file.jpg",
			want:  "albumName",
		},
		{
			name:  "album set an album's name based on folder path",
			album: "auto:folderPath",
			in:    "/foo/bar/file.jpg",
			want:  "foo_bar",
		},
		{
			name:  "album set an album's name based on folder name",
			album: "auto:folderName",
			in:    "/foo/bar/file.jpg",
			want:  "bar",
		},
		{
			name:  "album set an album's name with unexpected key (not `name` or `auto`)",
			album: "foo:bar",
			in:    "/foo/bar/file.jpg",
			want:  "",
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			job := UploadFolderJob{
				Album: tt.album,
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
		Album: "auto:fooBar",
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
