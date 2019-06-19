package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	// Version is this application version and it's set on build time
	Version = "v0.0.0"
	// Build is this specific build version and it's set on build time
	Build   = "0"
)

// versionCmd get the application version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current version and build",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gphotos-uploader-cli %s (build: %s)\n", Version, Build)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
