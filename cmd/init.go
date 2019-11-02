package cmd

import (
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"

	"github.com/gphotosuploader/gphotos-uploader-cli/cmd/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/log"
)

// InitCmd holds the required data for the init cmd
type InitCmd struct {
	CfgDir string
	// Flags
	Reconfigure bool
}

func NewInitCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &InitCmd{
		CfgDir: globalFlags.CfgDir,
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initializes configuration file",
		Long:  `Initializes a new configuration file. Creates a config.hjson with a default configuration.`,
		Args:  cobra.NoArgs,
		RunE:  cmd.Run,
	}

	initCmd.Flags().BoolVar(&cmd.Reconfigure, "force", false, "Overwrite existing configuration")

	return initCmd
}

func (cmd *InitCmd) Run(cobraCmd *cobra.Command, args []string) error {
	// Check if config already exists
	configExists := config.ConfigExists(cmd.CfgDir)
	if configExists && cmd.Reconfigure == false {
		log.Infof("Config already exists. If you want to recreate the config please run `%s`", ansi.Color("gphotos-uploader-cli init --force", "white+b"))
		log.Infof("\r         \nIf you want to continue with the existing config, run:\n- `%s` to start uploading files.\n", ansi.Color("gphotos-uploader-cli push", "white+b"))
		return nil
	}

	err := config.InitConfigFile(cmd.CfgDir)
	if err != nil {
		return err
	}

	log.Done("Configuration file successfully initialized.")
	log.Infof("\r         \nPlease edit: \n- `%s/config.hjson` to add you configuration.\n", cmd.CfgDir)

	return nil
}
