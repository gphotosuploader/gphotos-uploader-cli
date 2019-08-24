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
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			force = false
		}
		err = config.InitConfig(cfgDir, force)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Configuration file has been created.\nEdit it by running:\n    nano %s/config.hjson\n", cfgDir)
	},
}

func init() {
	initCmd.Flags().BoolP("force", "f", false, "Remove existing configuration folder (if exists)")
	rootCmd.AddCommand(initCmd)
}
