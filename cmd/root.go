package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/gphotosuploader/gphotos-uploader-cli/cmd/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/log"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "gphotos-uploader-cli",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Welcome to `gphotos-uploader-cli` a Google Photos uploader!",
	PersistentPreRun: func(cobraCmd *cobra.Command, args []string) {
		if globalFlags.Silent {
			log.GetInstance().SetLevel(logrus.FatalLevel)
		}
	},
	Long: `This application allows you to upload your pictures and videos to Google Photos. You can upload folders to several Google Photos accounts and organize them in albums.

    Get started by running the init command to configure your settings:
    $ gphotos-uploader-cli init

    once it's configured, start uploading your files:
    $ gphotos-uploader-cli push`,
}

var globalFlags *flags.GlobalFlags

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Execute command
	err := rootCmd.Execute()
	if err != nil {
		if globalFlags.Debug {
			log.Fatalf("%+v", err)
		} else {
			log.Fatal(err)
		}
		os.Exit(1)
	}
}

func init() {
	persistentFlags := rootCmd.PersistentFlags()
	globalFlags = flags.SetGlobalFlags(persistentFlags)

	// Add main commands
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewInitCmd(globalFlags))
	rootCmd.AddCommand(NewPushCmd(globalFlags))
}

// GetRoot returns the root command
func GetRoot() *cobra.Command {
	return rootCmd
}
