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
		isErrExpected bool
	}{
		{"Should success", "", "/home/foo/SourceFolder.hjson", false},
		{"Should success w/ existing dir", "/home/bar/SourceFolder.hjson", "/home/bar/SourceFolder.hjson", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			createTestConfigurationFile(t, fs, tc.preCreate)

			_, err := config.Create(fs, tc.path)
			assertExpectedError(t, tc.isErrExpected, err)

			if !tc.isErrExpected {
				assertFileExistence(t, fs, tc.path)
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
		{"Should return true if exist", "testdata/valid-config/config.hjson", true},
		{"Should return false if not exist", "testdata/non-existent/config.hjson", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.OsFs{}
			got := config.Exists(fs, tc.path)
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
		{"Should success with Album's name option", "testdata/valid-config/configWithAlbumNameOption.hjson", "youremail@domain.com", false},
		{"Should success with Album's auto folderName option", "testdata/valid-config/configWithAlbumAutoFolderNameOption.hjson", "youremail@domain.com", false},
		{"Should success with Album's auto folderPath option", "testdata/valid-config/configWithAlbumAutoFolderPathOption.hjson", "youremail@domain.com", false},
		{"Should success with deprecated CreateAlbums option", "testdata/valid-config/configWithDeprecatedCreateAlbumsOption.hjson", "youremail@domain.com", false},

		{"Should fail if config dir does not exist", "testdata/non-existent/config.hjson", "", true},
		{"Should fail if Account is invalid", "testdata/invalid-config/EmptyAccount.hjson", "", true},
		{"Should fail if SourceFolder does not exist", "testdata/invalid-config/NonExistentSourceFolder.hjson", "", true},
		{"Should fail if SecretsBackendType is invalid", "testdata/invalid-config/BadSecretsBackendType.hjson", "", true},
		{"Should fail if AppAPICredentials are invalid", "testdata/invalid-config/EmptyAppAPICredentials.hjson", "", true},
		{"Should fail if Jobs is empty", "testdata/invalid-config/NoJobs.hjson", "", true},
		{"Should fail if Album's format is invalid", "testdata/invalid-config/AlbumBadFormat.hjson", "", true},
		{"Should fail if Album's key is invalid", "testdata/invalid-config/AlbumBadKey.hjson", "", true},
		{"Should fail if Album's name is invalid", "testdata/invalid-config/AlbumEmptyName.hjson", "", true},
		{"Should fail if Album's auto value is invalid", "testdata/invalid-config/AlbumBadAutoValue.hjson", "", true},
		{"Should fail if deprecated CreateAlbums is invalid", "testdata/invalid-config/DeprecatedCreateAlbums.hjson", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.OsFs{}
			got, err := config.FromFile(fs, tc.path)
			if err != nil {
				t.Log(err)
			}
			assertExpectedError(t, tc.isErrExpected, err)

			if !tc.isErrExpected && (got.Account != tc.want) {
				t.Errorf("want: %s, got: %s", tc.want, got.Account)
			}
		})
	}
}
