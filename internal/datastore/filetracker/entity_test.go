package filetracker_test

import (
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/filetracker"
	"time"
)

func TestTrackedFile_Hash(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{"Should return the value", "123456789", "123456789"},
		{"Should return the value when it is in old format", "x|123456789", "123456789"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := filetracker.NewTrackedFile(tc.input)
			got := f.Hash
			if tc.want != got {
				t.Errorf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}

func TestTrackedFile_ModTime(t *testing.T) {
	testCases := []struct {
		name string
		input string
		want time.Time
	}{
		{"Should return zero time value", "123456789", time.Time{}},
		{"Should return time value", "1631350013816466000|123456789", time.Unix(0, 1631350013816466000)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := filetracker.NewTrackedFile(tc.input)
			got := f.ModTime
			if tc.want != got {
				t.Errorf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}

func TestTrackedFile_String(t *testing.T) {
	testCases := []struct {
		name string
		input string
		want string
	}{
		{"Should return the hash", "123456789", "123456789"},
		{"Should return mtime and hash", "1631350013816466000|123456789", "1631350013816466000|123456789"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f := filetracker.NewTrackedFile(tc.input)
			got := f.String()
			if tc.want != got {
				t.Errorf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
