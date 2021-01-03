package filetracker_test

import (
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/filetracker"
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
			got := f.Hash()
			if tc.want != got {
				t.Errorf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
