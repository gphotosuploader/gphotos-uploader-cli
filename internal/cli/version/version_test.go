package version_test

import (
	"bytes"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/configuration"
	versioninfo "github.com/gphotosuploader/gphotos-uploader-cli/version"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCommand(t *testing.T) {
	// Prepare a known version without depending on the build info.
	versioninfo.VersionInfo = &versioninfo.Info{
		Application:   "fooBarCommand",
		VersionString: "fooBarVersion",
	}

	actual := new(bytes.Buffer)

	configuration.Settings = configuration.Init("")
	rootCommand := cli.NewCommand()
	rootCommand.SetOut(actual)
	rootCommand.SetArgs([]string{"version"})

	_ = rootCommand.Execute()

	expected := "fooBarCommand Version: fooBarVersion\n"

	assert.Equal(t, expected, actual.String())
}
