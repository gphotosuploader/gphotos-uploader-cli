package configuration

import (
	"fmt"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/feedback"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

// Settings is a global instance of viper holding configurations for the CLI.
var Settings *viper.Viper

// Init initialize defaults and read the configuration file.
// Please note the logging system hasn't been configured yet,
// so logging shouldn't be used here.
func Init(configFile string) *viper.Viper {

	// Create a new viper instance with default values for all the settings
	settings := viper.New()
	SetDefaults(settings)

	// Set config name and config path
	if configFile != "" {
		settings.SetConfigName(strings.TrimSuffix(filepath.Base(configFile), filepath.Ext(configFile)))
		settings.AddConfigPath(filepath.Dir(configFile))
	} else {
		configDir := settings.GetString("directories.Data")
		// Get the default data path if none was provided
		if configDir == "" {
			configDir = getDefaultGooglePhotosCLIDataDir()
		}

		settings.SetConfigName("config")
		settings.AddConfigPath(configDir)
	}

	// Attempt to read config file
	if err := settings.ReadInConfig(); err != nil {
		// ConfigFileNotFoundError is acceptable, anything else
		// should be reported to the user
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			feedback.Warning(fmt.Sprintf("Error reading config file: %v", err))
		}
	}

	return settings
}

// BindFlags creates all the flags binding between the cobra Command and the instance of viper.
func BindFlags(cmd *cobra.Command, settings *viper.Viper) {
	_ = settings.BindPFlag("logging.level", cmd.Flag("log-level"))
	_ = settings.BindPFlag("output.no_color", cmd.Flag("no-color"))
}

// getDefaultGooglePhotosCLIDataDir returns the full path to the default gphotos-cli folder.
func getDefaultGooglePhotosCLIDataDir() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		feedback.Warning(fmt.Sprintf("Unable to get user home dir: %v", err))
		return "."
	}

	return filepath.Join(userHomeDir, ".gphotos-cli")
}

// FindConfigFileInArgs returns the config file path using the
// argument '--config-file' (if specified) or looking in the current working dir.
func FindConfigFileInArgs(args []string) string {
	// Look for '--config-file' argument
	for i, arg := range args {
		if arg == "--config-file" {
			if len(args) > i+1 {
				return args[i+1]
			}
		}
	}
	return ""
}
