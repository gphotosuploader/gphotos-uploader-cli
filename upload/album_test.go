package upload

import (
	"testing"
)

func TestAlbumNameUsingFolderPath(t *testing.T) {
	var testData = []struct {
		in   string
		out  string
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
		in   string
		out  string
	}{
		{in: "", out: ""},
		{in: "foo", out: ""},
		{in: "foo/", out: "foo"},
		{in: "foo/bar", out: "bar"},
		{in: "foo/bar/", out: "bar"},
	}
	for _, tt := range testData {
		got := albumNameUsingFolderName(tt.in)
		if got != tt.out {
			t.Errorf("albumNameUsingFolderName for '%s' failed: expected '%s', got '%s'", tt.in, tt.out, got)
		}
	}
}

