package version_test

import (
	"bytes"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/version"
	versioninfo "github.com/gphotosuploader/gphotos-uploader-cli/version"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCommand(t *testing.T) {
	// Prepare a know version without depending on the build info.
	versioninfo.VersionInfo = &versioninfo.Info{
		Application:   "fooBarCommand",
		VersionString: "fooBarVersion",
	}

	actual := new(bytes.Buffer)
	versionCommand := version.NewCommand()
	versionCommand.SetOut(actual)
	_ = versionCommand.Execute()

	expected := "fooBarCommand Version: fooBarVersion\n"

	assert.Equal(t, expected, actual.String())
}
