package cli

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/auth"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/list"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/push"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/version"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/configuration"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/feedback"
	versioninfo "github.com/gphotosuploader/gphotos-uploader-cli/version"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

var (
	longCommandDescription = `
Google Photos Command Line Interface (CLI)

This CLI application allows you to upload pictures and videos to Google Photos. You can upload folders to your Google Photos account and organize them in albums automatically. Additionally, you can list albums and media items already uploaded to Google Photos.

To get started, initialize your settings by running the following command:
$ gphotos-cli init

Once configured, you can uploading your files with this command:
$ gphotos-cli push

Or you can list your albums in Google Photos by running:
$ gphotos-cli list albums

For more information, visit: https://gphotosuploader.github.io/gphotos-uploader-cli.
`

	verbose    bool
	configFile string

	// Os points to the (real) file system.
	// Useful for testing.
	Os = afero.NewOsFs()
)

// NewCommand creates a new gphotosCLI command root
func NewCommand() *cobra.Command {
	// gphotosCLI is the root command
	gphotosCLI := &cobra.Command{
		Use:              "gphotos-cli",
		Short:            "Google Photos CLI.",
		Long:             longCommandDescription,
		PersistentPreRun: preRun,
	}

	createCliCommandTree(gphotosCLI)

	return gphotosCLI
}

// this is here only for testing
func createCliCommandTree(cmd *cobra.Command) {
	// Add main commands
	cmd.AddCommand(version.NewCommand())
	cmd.AddCommand(config.NewCommand())
	cmd.AddCommand(push.NewCommand())
	cmd.AddCommand(auth.NewCommand())
	cmd.AddCommand(list.NewCommand())

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Print the logs on the standard output.")

	validLogLevels := []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}
	cmd.PersistentFlags().String("log-level", "", fmt.Sprintf("Messages with this level and above will be logged. Valid levels are: %s", strings.Join(validLogLevels, ", ")))

	cmd.PersistentFlags().StringVar(&configFile, "config-file", "", "The custom config file (if not specified the default will be used).")
	cmd.PersistentFlags().Bool("no-color", false, "Disable colored output.")
	configuration.BindFlags(cmd, configuration.Settings)
}

func preRun(cobraCmd *cobra.Command, args []string) {
	configFile := configuration.Settings.ConfigFileUsed()

	// https://no-color.org/
	color.NoColor = configuration.Settings.GetBool("output.no_color") || os.Getenv("NO_COLOR") != ""

	// Set default feedback output to colorable
	feedback.SetOut(colorable.NewColorableStdout())
	feedback.SetErr(colorable.NewColorableStderr())

	//
	// Prepare logging
	//

	// decide whether we should log to stdout
	if verbose {
		// if we print on stdout, do it in full colors
		logrus.SetOutput(colorable.NewColorableStdout())
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			DisableColors: color.NoColor,
		})
	} else {
		logrus.SetOutput(io.Discard)
	}

	// configure logging filter
	if lvl, found := toLogLevel(configuration.Settings.GetString("logging.level")); !found {
		feedback.Fatal(fmt.Sprintf("Invalid option for --log-level: %s", configuration.Settings.GetString("logging.level")), feedback.ErrBadArgument)
	} else {
		logrus.SetLevel(lvl)
	}

	//
	// Print some status info and check command is consistent
	//

	if configFile != "" {
		logrus.Infof("Using config file: %s", configFile)
	} else {
		logrus.Info("Config file not found, using default values")
	}

	logrus.Info(versioninfo.VersionInfo.Application + " version " + versioninfo.VersionInfo.VersionString)
}

// convert the string passed to the `--log-level` option to the corresponding
// logrus formal level.
func toLogLevel(s string) (t logrus.Level, found bool) {
	t, found = map[string]logrus.Level{
		"trace": logrus.TraceLevel,
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"panic": logrus.PanicLevel,
	}[s]

	return
}
