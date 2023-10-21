package cli_test

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/flags"
)

func TestNewInitCmd(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		args          []string
		isErrExpected bool
	}{
		{"Should success", "", []string{}, false},
		{"Should fail if input exists", "/foo", []string{}, true},
		{"Should success if input exists and force is set", "/foo", []string{"--force"}, false},
	}

	t.Cleanup(func() {
		cli.Os = afero.NewOsFs()
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cli.Os = afero.NewMemMapFs()
			createTestConfigurationFile(t, cli.Os, tc.input)

			c := cli.NewInitCmd(&flags.GlobalFlags{CfgDir: tc.input})
			c.SetArgs(tc.args)

			err := c.Execute()
			assertExpectedError(t, tc.isErrExpected, err)
		})
	}
}

func createTestConfigurationFile(t *testing.T, fs afero.Fs, path string) {
	if path == "" {
		return
	}
	if err := fs.MkdirAll(path, 0700); err != nil {
		t.Fatalf("creating test dir, err: %s", err)
	}
	filename := filepath.Join(path, app.DefaultConfigFilename)
	if err := afero.WriteFile(fs, filename, []byte("my"), 0600); err != nil {
		t.Fatalf("creating test configuration file, err: %s", err)
	}
}
