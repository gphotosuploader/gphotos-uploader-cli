package cmd

import (
	"path/filepath"

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
		Short: "Initializes the configuration",
		Long:  `Initializes a new configuration with defaults.`,
		Args:  cobra.NoArgs,
		RunE:  cmd.Run,
	}

	initCmd.Flags().BoolVar(&cmd.Reconfigure, "force", false, "Overwrite existing configuration")

	return initCmd
}

func (cmd *InitCmd) Run(cobraCmd *cobra.Command, args []string) error {
	path, _ := filepath.Abs(cmd.CfgDir)

	// Check if config already exists
	if config.Exists(path) && !cmd.Reconfigure {
		log.Infof("Configuration file already exists. If you proceed, %s", ansi.Color("ALL THE APPLICATION DATA WILL BE DELETED!", "white+b"))
		log.Infof("Use `%s` flag to proceed and recreate the configuration file", ansi.Color("--force", "white+b"))
		return nil
	}

	file, err := recreateAppDir(path)
	if err != nil {
		log.Failf("Unable to create configuration file, err: %s", err)
		return err
	}

	log.Done("Configuration file created successfully.")
	log.Infof("\r         \nPlease edit: \n- `%s` to add you configuration.\n",
		ansi.Color(file, "cyan+b"),
	)

	return nil
}

// recreateAppDir returns the configuration file name created after creating the application directory.
func recreateAppDir(path string) (string, error) {
	if err := emptyDir(path); err != nil {
		return "", err
	}
	return config.Create(path)
}

func emptyDir(path string) error {
	if err := Os.RemoveAll(path); err != nil {
		return err
	}
	return Os.MkdirAll(path, 0700)
}
