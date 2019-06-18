package cmd

import (
	"fmt"
	"github.com/nmrshll/gphotos-uploader-cli/config"
	"github.com/spf13/cobra"
	"log"
)

// versionCmd get the application version
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init the configuration file",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.InitConfigFile(defaultCfgFile)
		if err != nil {
			log.Fatalf("Failed to create the init config file: %v", err)
		}
		fmt.Printf("Configuration file has been created.\nEdit it by running:\n    nano %s\n", defaultCfgFile)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
