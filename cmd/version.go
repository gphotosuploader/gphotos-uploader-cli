package cmd

import (
	"fmt"
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
	version = "v0.0.0" // git tag, output of $(git describe --tags --abbrev=0)
	build   = "0"      // sha1 from git, output of $(git rev-parse --short HEAD)
)

// versionCmd returns the overall codebase version. It's for detecting what
// code a binary was built from.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current version and build",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gphotos-uploader-cli %s (build: %s)\n", version, build)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
