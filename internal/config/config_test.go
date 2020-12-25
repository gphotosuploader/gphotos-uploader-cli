package config_test

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		name          string
		preCreate     string
		path          string
		want          string
		isErrExpected bool
	}{
		{"Should success", "", "/home/foo", "/home/foo/config.hjson", false},
		{"Should success w/ existing dir", "/home/bar", "/home/bar", "/home/bar/config.hjson", false},
	}

	t.Cleanup(func() {
		config.Os = afero.NewOsFs()
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config.Os = afero.NewMemMapFs()
			createTestDir(t, tc.preCreate)

			got, err := config.Create(tc.path)
			assertExpectedError(t, tc.isErrExpected, err)
			if !tc.isErrExpected && tc.want != got {
				t.Errorf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}

func TestExists(t *testing.T) {
	testCases := []struct {
		name string
		path string
		want bool
	}{
		{"Should return true if exist", "testdata/valid-config", true},
		{"Should return false if not exist", "testdata/non-existent", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := config.Exists(tc.path)
			if tc.want != got {
				t.Errorf("configuration file does not exist, path: %s", tc.path)
			}
		})
	}
}

func TestFromFile(t *testing.T) {
	testCases := []struct {
		name          string
		path          string
		want          string
		isErrExpected bool
	}{
		{"Should success", "testdata/valid-config", "youremail@domain.com", false},
		{"Should fail if dir does not exist", "testdata/non-existent", "", true},
		{"Should fail if config data is invalid", "testdata/invalid-config", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := config.FromFile(tc.path)
			assertExpectedError(t, tc.isErrExpected, err)

			if !tc.isErrExpected && (got.Account != tc.want) {
				t.Errorf("want: %s, got: %s", tc.want, got.Account)
			}
		})
	}
}

func createTestDir(t *testing.T, path string) {
	if path == "" {
		return
	}
	if err := config.Os.MkdirAll(path, 0700); err != nil {
		t.Fatalf("error creating test dir, err: %s", err)
	}
}

func assertExpectedError(t *testing.T, errExpected bool, err error) {
	if errExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !errExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}
