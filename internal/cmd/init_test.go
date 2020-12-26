package cmd_test

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
)

func TestNewInitCmd(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		args          []string
		force         bool
		isErrExpected bool
	}{
		{"Should success", "", []string{}, false, false},
		{"Should fail if input exists", "/foo", []string{}, false, false},
		{"Should success if input exists and force is set", "/foo", []string{"--force"}, false, false},
	}

	t.Cleanup(func() {
		cmd.Os = afero.NewOsFs()
		config.Os = cmd.Os
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd.Os = afero.NewMemMapFs()
			config.Os = cmd.Os
			createTestConfigurationFile(t, cmd.Os, tc.input)

			c := cmd.NewInitCmd(&flags.GlobalFlags{CfgDir: tc.input})
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
	filename := filepath.Join(path, config.DefaultConfigFilename)
	if err := afero.WriteFile(fs, filename, []byte("my"), 0600); err != nil {
		t.Fatalf("creating test configuration file, err: %s", err)
	}
}
