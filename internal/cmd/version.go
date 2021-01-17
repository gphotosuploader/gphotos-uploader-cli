package cmd

import (
	"github.com/spf13/cobra"
)

// Base version information.
//
// This is the fallback data used when version information from git is not
// provided via go ldflags. It provides an approximation of the application
// version for ad-hoc builds (e.g. `go build`) that cannot get the version
// information from git
//
// If you are looking at these fields in the git tree, they look strange. They
// are modified on the fly by the build process.
//
// We use semantic version (see https://semver.org/ for more information). When
// releasing a new version, this file is updated by Makefile to reflect the new
// version, a git annotated tag is used to set this version
var (
	version = "v0.0.0" // git tag, output of $(git describe --tags --always --dirty)
)

type VersionCmd struct{}

func NewVersionCmd() *cobra.Command {
	cmd := &VersionCmd{}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Prints current version",
		Run:  cmd.Run,
	}

	return versionCmd
}

func (cmd *VersionCmd) Run(command *cobra.Command, args []string) {
	command.Printf("gphotos-uploader-cli %s\n", version)
}
