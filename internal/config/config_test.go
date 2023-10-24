package config_test

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
)

func TestConfig_SafePrint(t *testing.T) {
	cfg := config.Config{
		APIAppCredentials: config.APIAppCredentials{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		},
		Account:            "account",
		SecretsBackendType: "auto",
		Jobs: []config.FolderUploadJob{
			{
				SourceFolder:      "foo",
				Album:             "name:albumName",
				CreateAlbums:      "folderPath",
				DeleteAfterUpload: false,
				IncludePatterns:   []string{},
				ExcludePatterns:   []string{},
			},
		},
	}
	want := `{"APIAppCredentials":{"ClientID":"client-id","ClientSecret":"REMOVED"},"Account":"account","SecretsBackendType":"auto","Jobs":[{"SourceFolder":"foo","Album":"name:albumName","CreateAlbums":"folderPath","DeleteAfterUpload":false,"IncludePatterns":[],"ExcludePatterns":[]}]}`

	if want != cfg.SafePrint() {
		t.Errorf("want: %s, got: %s", want, cfg.SafePrint())
	}
}

func createTestConfigurationFile(t *testing.T, fs afero.Fs, path string) {
	if path == "" {
		return
	}
	if err := fs.MkdirAll(filepath.Dir(path), 0700); err != nil {
		t.Fatalf("creating test dir, err: %s", err)
	}
	if err := afero.WriteFile(fs, path, []byte("my"), 0600); err != nil {
		t.Fatalf("creating test configuration file, err: %s", err)
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

func assertFileExistence(t *testing.T, fs afero.Fs, path string) {
	exist, err := afero.Exists(fs, path)
	if err != nil {
		t.Fatalf("checking file existence, err: %s", err)
	}
	if !exist {
		t.Errorf("file expected, but it does not exist")
	}
}
