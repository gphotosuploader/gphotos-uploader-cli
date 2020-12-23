package cmd

import (
	"fmt"

	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

// InitCmd holds the required data for the init cmd
type InitCmd struct {
	*flags.GlobalFlags

	// command flags
	Reconfigure bool
}

func NewInitCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &InitCmd{GlobalFlags: globalFlags}

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
	if config.Exists(cmd.CfgDir) && !cmd.Reconfigure {
		log.Infof("AppConfig already exists. If you want to recreate the config please run `%s`", ansi.Color("gphotos-uploader-cli init --force", "white+b"))
		log.Infof("\r         \nIf you want to continue with the existing config, run:\n- `%s` to start uploading files.\n", ansi.Color("gphotos-uploader-cli push", "white+b"))
		return nil
	}

	if err := config.InitConfigFile(cmd.CfgDir); err != nil {
		return err
	}

	log.Done("Configuration file successfully initialized.")
	log.Infof("\r         \nPlease edit: \n- `%s` to add you configuration.\n",
		ansi.Color(fmt.Sprintf("%s/%s", cmd.CfgDir, config.DefaultConfigFilename), "cyan+b"),
	)

	return nil
}
