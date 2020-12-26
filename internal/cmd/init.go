package cmd

import (
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd/flags"
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
		Short: "Initializes the configuration",
		Long:  `Initializes a new configuration with defaults.`,
		Args:  cobra.NoArgs,
		RunE:  cmd.Run,
	}

	initCmd.Flags().BoolVar(&cmd.Reconfigure, "force", false, "Overwrite existing configuration")

	return initCmd
}

func (cmd *InitCmd) Run(cobraCmd *cobra.Command, args []string) error {
	cli, err := app.StartWithoutConfig(Os, cmd.CfgDir)
	if err != nil {
		return err
	}

	if exist := cli.AppDataDirExists(); exist && !cmd.Reconfigure {
		log.Infof("Application data already exists. If you proceed, %s", ansi.Color("ALL THE APPLICATION DATA WILL BE DELETED!", "white+b"))
		log.Infof("Use `%s` flag to proceed and recreate the application data", ansi.Color("--force", "white+b"))
		return nil
	}

	filename, err := cli.CreateAppDataDir()
	if err != nil {
		log.Failf("Unable to create application data dir, err: %s", err)
		return err
	}

	log.Done("Application data dir created successfully.")
	log.Infof("\r         \nPlease edit: \n- `%s` to add you configuration.\n",
		ansi.Color(filename, "cyan+b"),
	)

	return nil
}

