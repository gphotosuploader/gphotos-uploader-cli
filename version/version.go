package version

import (
	"fmt"
)

// VersionInfo contains all info injected during build
var VersionInfo *Info

// Base version information.
// We use semantic version (see https://semver.org/ for more information).
var (
	// When releasing a new version, Makefile updates the versionString to reflect the new
	// version; a git-annotated tag is used to set this version.
	versionString = "" // git tag, output of $(git describe --tags --always --dirty)

	// This is the fallback data used when version information from git is not
	// provided via go ldflags. It provides an approximation of the application
	// version for adhoc builds (e.g. `go build`) that cannot get the version
	// information from git
	defaultVersionString = "0.0.0-git"
)

type Info struct {
	Application   string `json:"Application"`
	VersionString string `json:"VersionString"`
}

func NewInfo() *Info {
	return &Info{
		Application:   "gphotos-uploader-cli",
		VersionString: versionString,
	}
}

func (i *Info) String() string {
	return fmt.Sprintf("%s Version: %s", i.Application, i.VersionString)
}

func init() {
	if versionString == "" {
		versionString = defaultVersionString
	}

	VersionInfo = NewInfo()
}
