package filetracker_test

import (
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/filetracker"
	"testing"
)

func TestXXHash32Hasher_Hash(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		want          string
		isErrExpected bool
	}{
		{"Should success", "testdata/image.jpg", "1127908779", false},
		{"Should fail", "testdata/non-existent", "", true},
	}

	hasher := filetracker.XXHash32Hasher{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := hasher.Hash(tc.input)
			assertExpectedError(t, tc.isErrExpected, err)
			if tc.want != got {
				t.Errorf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
