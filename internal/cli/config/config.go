package config

import (
	"github.com/spf13/cobra"
)

// NewCommand created a new `config` command
func NewCommand() *cobra.Command {
	configCommand := &cobra.Command{
		Use:   "config",
		Short: "Configuration commands.",
	}

	configCommand.AddCommand(initDumpCommand())
	configCommand.AddCommand(initInitCommand())

	return configCommand
}
