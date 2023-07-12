package flags

import (
	"github.com/mitchellh/go-homedir"
	flag "github.com/spf13/pflag"
)

const (
	// defaultApplicationDataFolder is the default folder where to store application data.
	defaultApplicationDataFolder = "~/.gphotos-cli"
)

// GlobalFlags is the flags that contain the global flags
type GlobalFlags struct {
	Silent bool
	Debug  bool
	CfgDir string
}

// SetGlobalFlags applies the global flags
func SetGlobalFlags(flags *flag.FlagSet) *GlobalFlags {
	globalFlags := &GlobalFlags{}

	flags.BoolVar(&globalFlags.Debug, "debug", false, "Log very verbose information. Useful for troubleshooting.")
	flags.BoolVar(&globalFlags.Silent, "silent", false, "Run in silent mode and prevent any log output except panics.")

	flags.StringVar(&globalFlags.CfgDir, "config", defaultApplicationDataPath(), "Sets the config folder path. All configuration will be kept in this folder.")

	return globalFlags
}

// defaultApplicationDataPath is the path to the default folder where to store application data.
func defaultApplicationDataPath() string {
	absPath, err := homedir.Expand(defaultApplicationDataFolder)
	if err != nil {
		return defaultApplicationDataFolder
	}
	return absPath
}
