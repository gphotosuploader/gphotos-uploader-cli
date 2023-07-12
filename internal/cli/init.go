package cli

import (
	"fmt"

	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

// InitCmd holds the required data for the init cmd
type InitCmd struct {
	*flags.GlobalFlags

	// command flags
	Force bool
}

func NewInitCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &InitCmd{GlobalFlags: globalFlags}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the configuration",
		Long:  `Initialize a new configuration with defaults.`,
		Args:  cobra.NoArgs,
		RunE:  cmd.Run,
	}

	initCmd.Flags().BoolVar(&cmd.Force, "force", false, "Overwrite existing configuration")

	return initCmd
}

func (cmd *InitCmd) Run(cobraCmd *cobra.Command, args []string) error {
	cli, err := app.StartWithoutConfig(Os, cmd.CfgDir)
	if err != nil {
		return err
	}

	if exist := cli.AppDataDirExists(); exist && !cmd.Force {
		log.Infof("Application data already exists. Use `%s` flag to proceed. %s",
			ansi.Color("--force", "white+b"),
			ansi.Color("ALL THE APPLICATION DATA WILL BE DELETED!", "white+b"))
		return fmt.Errorf("application data already exists at %s", cmd.CfgDir)
	}

	filename, err := cli.CreateAppDataDir()
	if err != nil {
		log.Failf("Unable to create application data dir, err: %s", err)
		return err
	}

	log.Done("Application data dir created successfully.")
	log.Infof("\r         \nPlease edit: \n- `%s` to add your configuration.\n",
		ansi.Color(filename, "cyan+b"),
	)

	return nil
}
