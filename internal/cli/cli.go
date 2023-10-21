package cli

import (
	"fmt"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/auth"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/list"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/push"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/version"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

var (
	longCommandDescription = `
Google Photos Command Line Interface (CLI)

This CLI application allows you to upload pictures and videos to Google Photos. You can upload folders to your Google Photos account and organize them in albums automatically. Additionally, you can list albums and media items already uploaded to Google Photos.

To get started, initialize your settings by running the following command:
$ gphotos-uploader-cli init

Once configured, you can uploading your files with this command:
$ gphotos-uploder-cli push

Or you can list your albums in Google Photos by running:
$ gphotos-uploader-cli list albums

For more information, visit: https://gphotosuploader.github.io/gphotos-uploader-cli.
`
	globalFlags *flags.GlobalFlags

	// Os points to the (real) file system.
	// Useful for testing.
	Os = afero.NewOsFs()
)

// NewCommand creates a new gphotosCLI command root
func NewCommand() *cobra.Command {
	// ArduinoCli is the root command
	gphotosCLI := &cobra.Command{
		Use:               "gphotos-uploader-cli",
		Short:             "Google Photos CLI.",
		Long:              longCommandDescription,
		PersistentPreRunE: preRun,
	}

	createCliCommandTree(gphotosCLI)

	return gphotosCLI
}

// this is here only for testing
func createCliCommandTree(cmd *cobra.Command) {
	persistentFlags := cmd.PersistentFlags()
	globalFlags = flags.SetGlobalFlags(persistentFlags)

	// Add main commands
	cmd.AddCommand(version.NewCommand())
	cmd.AddCommand(NewInitCmd(globalFlags))
	cmd.AddCommand(push.NewCommand(globalFlags))
	cmd.AddCommand(auth.NewCommand(globalFlags))
	cmd.AddCommand(list.NewCommand(globalFlags))

	// TODO: Set flags here instead of passing globalFlags to all commands.
	// See: https://github.com/arduino/arduino-cli/blob/master/internal/cli/cli.go
	//
	//cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print the logs on the standard output.")
	//validLogLevels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	//cmd.PersistentFlags().String("log-level", "", fmt.Sprintf("Messages with this level and above will be logged. Valid levels are: %s", strings.Join(validLogLevels, ", ")))
	//cmd.RegisterFlagCompletionFunc("log-level", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	//    return validLogLevels, cobra.ShellCompDirectiveDefault
	//})
	//cmd.PersistentFlags().String("log-file", "", "Path to the file where logs will be written.")
	//validLogFormats := []string{"text", "json"}
	//cmd.PersistentFlags().String("log-format", "", fmt.Sprintf("The output format for the logs, can be: %s", strings.Join(validLogFormats, ", ")))
	//cmd.RegisterFlagCompletionFunc("log-format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	//    return validLogFormats, cobra.ShellCompDirectiveDefault
	//})
	//cmd.PersistentFlags().StringVar(&configFile, "config-file", "", "The custom config file (if not specified the default will be used).")
	//cmd.PersistentFlags().Bool("no-color", false, "Disable colored output.")
}

func preRun(cobraCmd *cobra.Command, args []string) error {
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
}
