package flags

import (
	"github.com/mitchellh/go-homedir"
	flag "github.com/spf13/pflag"
)

const (
	// defaultApplicationDataFolder is the default folder where to store application data.
	defaultApplicationDataFolder = "~/.gphotos-uploader-cli"
)

// GlobalFlags is the flags that contains the global flags
type GlobalFlags struct {
	Silent bool
	Debug  bool
	CfgDir string
}

// SetGlobalFlags applies the global flags
func SetGlobalFlags(flags *flag.FlagSet) *GlobalFlags {
	globalFlags := &GlobalFlags{}

	flags.BoolVar(&globalFlags.Debug, "debug", false, "Logs very verbose information. Useful for troubleshooting.")
	flags.BoolVar(&globalFlags.Silent, "silent", false, "Run in silent mode and prevents any log output except panics & fatals.")

	flags.StringVar(&globalFlags.CfgDir, "config", defaultApplicationDataPath(), "Sets config folder path. All configuration will be keep in this folder.")

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
