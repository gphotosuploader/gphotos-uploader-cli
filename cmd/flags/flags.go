package flags

import flag "github.com/spf13/pflag"

const defaultCfgDir = "~/.config/gphotos-uploader-cli"

// GlobalFlags is the flags that contains the global flags
type GlobalFlags struct {
	Silent bool
	Debug  bool
	CfgDir string
}

// SetGlobalFlags applies the global flags
func SetGlobalFlags(flags *flag.FlagSet) *GlobalFlags {
	globalFlags := &GlobalFlags{
		CfgDir: defaultCfgDir,
	}

	flags.BoolVar(&globalFlags.Debug, "debug", false, "Logs very verbose information. Useful for troubleshooting.")
	flags.BoolVar(&globalFlags.Silent, "silent", false, "Run in silent mode and prevents any log output except panics & fatals.")
	flags.StringVar(&globalFlags.CfgDir, "config", defaultCfgDir, "Sets config folder path. All configuration will be keep in this folder.")

	return globalFlags
}
