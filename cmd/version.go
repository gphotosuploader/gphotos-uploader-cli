package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	Version = "v0.0.0"
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