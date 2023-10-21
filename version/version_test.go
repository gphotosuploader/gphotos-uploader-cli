package version

import (
	"github.com/stretchr/testify/require"
	"testing"
)

// TestBuildInjectedInfo tests the Info strings passed to the binary at build time
// in order to have this test green launch your testing 'make test' or use:
//
//	go test -run TestBuildInjectedInfo -v ./... -ldflags '
//	  -X github.com/gphotosuploader/gphotos-uploader-cli/version.versionString=0.0.0-test.preview'
func TestBuildInjectedInfo(t *testing.T) {
	goldenInfo := Info{
		Application:   "gphotos-uploader-cli",
		VersionString: "0.0.0-test.preview",
	}
	info := NewInfo()
	require.Equal(t, goldenInfo.Application, info.Application)
	require.Equal(t, goldenInfo.VersionString, info.VersionString)
	require.Equal(t, "gphotos-uploader-cli Version: 0.0.0-test.preview", info.String())
}
