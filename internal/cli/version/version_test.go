package version_test

import (
	"bytes"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/configuration"
	versioninfo "github.com/gphotosuploader/gphotos-uploader-cli/version"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewVersionCommand(t *testing.T) {
	// Prepare a known version without depending on the build info.
	versioninfo.VersionInfo = &versioninfo.Info{
		Application:   "fooBarCommand",
		VersionString: "fooBarVersion",
	}
	expected := "fooBarCommand Version: fooBarVersion\n"

	actual := new(bytes.Buffer)
	configuration.Settings = configuration.Init("")
	rootCommand := cli.New()
	rootCommand.SetOut(actual)
	rootCommand.SetArgs([]string{"version"})

	err := rootCommand.Execute()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual.String())
}
