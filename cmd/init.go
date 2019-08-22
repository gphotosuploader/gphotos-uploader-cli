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
		err := config.InitConfig(defaultCfgDir)
		if err != nil {
			log.Fatalf("Failed to create the configuration: path=%s, err=%v", defaultCfgDir, err)
		}
		fmt.Printf("Configuration file has been created.\nEdit it by running:\n    nano %s/config.hjson\n", defaultCfgDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
