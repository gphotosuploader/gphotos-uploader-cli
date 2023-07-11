package version

import (
	"github.com/spf13/cobra"
)

// Base version information.
// We use semantic version (see https://semver.org/ for more information).
var (
	// When releasing a new version, Makefile updates this file to reflect the new
	// version; a git-annotated tag is used to set this version.
	version = "v0.0.0" // git tag, output of $(git describe --tags --always --dirty)

	// This is the fallback data used when version information from git is not
	// provided via go ldflags. It provides an approximation of the application
	// version for adhoc builds (e.g. `go build`) that cannot get the version
	// information from git
	defaultVersionString = "0.0.0-git"
)

func NewCommand() *cobra.Command {
	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "Shows version number of gphotos-cli.",
		Long:  "Shows the version number of Arduino CLI which is installed on your system.",
		Args:  cobra.NoArgs,
		Run:   runVersionCommand,
	}

	return versionCommand
}

func runVersionCommand(cmd *cobra.Command, args []string) {
	//if version == "" {
	//    version = defaultVersionString
	//}

	cmd.Printf("gphotos-cli %s\n", version)
}
