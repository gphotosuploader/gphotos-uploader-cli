package filetracker_test

import (
	"testing"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/filetracker"
)

func TestXxHash32Hasher_Hash(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		want          string
		isErrExpected bool
	}{
		{"Should success", "testdata/image.jpg", "1127908779", false},
		{"Should fail", "testdata/non-existent", "", true},
	}

	ft := filetracker.New(&mockedRepository{})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ft.Hasher.Hash(tc.input)
			assertExpectedError(t, tc.isErrExpected, err)
			if tc.want != got {
				t.Errorf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
