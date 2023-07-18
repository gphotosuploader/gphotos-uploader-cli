package configuration

import (
	"github.com/spf13/viper"
	"strings"
)

// SetDefaults sets the default values for certain keys
func SetDefaults(settings *viper.Viper) {
	// logging
	settings.SetDefault("logging.level", "info")

	// gphotos-cli directories
	settings.SetDefault("directories.data", getDefaultGooglePhotosCLIDataDir())

	// output settings
	settings.SetDefault("output.no_color", false)

	// auth settings
	settings.SetDefault("auth.secrets_type", "auto")

	// Bind env vars
	settings.SetEnvPrefix("GPHOTOS_CLI")
	settings.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	settings.AutomaticEnv()

	// Bind env aliases to keep backward compatibility
	_ = settings.BindEnv("directories.Data", "GPHOTOS_CLI_DATA_DIR")
}
