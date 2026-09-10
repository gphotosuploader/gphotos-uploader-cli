package filetracker_test

import (
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/filetracker"
	"testing"
)

func TestSHA256Hasher_Hash(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		want          string
		isErrExpected bool
	}{
		{"Should success", "testdata/image.jpg", "b554e7ea1a1485c86e6c97c387d4f0f13a08114502e71bddad5482e6fa53cbae", false},
		{"Should fail", "testdata/non-existent", "", true},
	}

	hasher := filetracker.SHA256Hasher{}

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
