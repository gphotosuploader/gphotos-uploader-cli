package upload

import (
	"testing"
)

func TestAlbumName(t *testing.T) {
	var testData = []struct {
		name               string
		createAlbum        bool
		createAlbumBasedOn string

		in   string
		want string
	}{
		{
			name:               "createAlbumDisabled_With_folderName",
			createAlbum:        false,
			createAlbumBasedOn: "folderName",
			in:                 "/foo/bar/file.jpg",
			want:               "",
		},
		{
			name:               "createAlbumDisabled_With_folderPath",
			createAlbum:        false,
			createAlbumBasedOn: "folderPath",
			in:                 "/foo/bar/file.jpg",
			want:               "",
		},
		{
			name:               "createAlbum_With_folderName",
			createAlbum:        true,
			createAlbumBasedOn: "folderName",
			in:                 "/foo/bar/file.jpg",
			want:               "bar",
		},
		{
			name:               "createAlbum_With_folderPath",
			createAlbum:        true,
			createAlbumBasedOn: "folderPath",
			in:                 "/foo/bar/file.jpg",
			want:               "foo_bar",
		},
	}

	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			job := UploadFolderJob{
				CreateAlbum:        tt.createAlbum,
				CreateAlbumBasedOn: tt.createAlbumBasedOn,
			}
			got := job.albumName(tt.in)
			if got != tt.want {
				t.Errorf("albumName for '%s' failed: expected '%s', got '%s'", tt.in, tt.want, got)
			}
		})

	}
}

func TestAlbumNameUsingTemplate(t *testing.T) {
	var testData = []struct {
		in       string
		template string
		out      string
	}{
		{in: "/foo/bar/xyz", template: "folderPath", out: "foo_bar"},
		{in: "foo/bar/xyz", template: "folderPath", out: "foo_bar"},
		{in: "foo/bar/xyz", template: "folderName", out: "bar"},
		{in: "/foo/bar/xyz/file", template: "folderPath", out: "foo_bar_xyz"},
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
