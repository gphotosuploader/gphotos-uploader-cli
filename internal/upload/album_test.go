package upload

import (
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
		{
			name:         "createAlbum_With_customName",
			createAlbums: "MyAlbum",
			in:           "/foo/bar/file.jpg",
			want:         "MyAlbum",
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			job := UploadFolderJob{
				CreateAlbums: tt.createAlbums,
			}
			got := job.albumName(tt.in)
			if got != tt.want {
				t.Errorf("albumName for '%s' failed: expected '%s', got '%s'", tt.in, tt.want, got)
			}
		})

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
		{in: "/foo/bar", out: "foo"},
		{in: "/foo/bar/", out: "foo_bar"},
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
		{in: "/foo/bar", out: "foo"},
		{in: "/foo/bar/", out: "bar"},
	}
	for _, tt := range testData {
		got := albumNameUsingFolderName(tt.in)
		if got != tt.out {
			t.Errorf("albumNameUsingFolderName for '%s' failed: expected '%s', got '%s'", tt.in, tt.out, got)
		}
	}
}
