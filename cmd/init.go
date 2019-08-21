package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/gphotosuploader/gphotos-uploader-cli/config"
)

// versionCmd get the application version
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.InitConfigFile(defaultCfgDir)
		if err != nil {
			log.Fatalf("Failed to create the init config file: %v", err)
		}
		fmt.Printf("Configuration file has been created.\nEdit it by running:\n    nano %s\n", defaultCfgDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
