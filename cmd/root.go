package cmd

import (
	"fmt"
	"os"

	"github.com/mgutz/ansi"
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
	PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
		if globalFlags.Silent && globalFlags.Debug {
			return fmt.Errorf("%s and %s cannot be specified at the same time", ansi.Color("--silent", "white+b"), ansi.Color("--debug", "white+b"))
		}
		if globalFlags.Silent {
			log.GetInstance().SetLevel(logrus.FatalLevel)
		}
		if globalFlags.Debug {
			log.GetInstance().SetLevel(logrus.DebugLevel)
		}
		return nil
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
	rootCmd.AddCommand(NewAuthCmd(globalFlags))
}

// GetRoot returns the root command
func GetRoot() *cobra.Command {
	return rootCmd
}
