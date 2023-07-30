package config

import (
    "fmt"
    "github.com/gphotosuploader/gphotos-uploader-cli/internal/configuration"
    "github.com/gphotosuploader/gphotos-uploader-cli/internal/feedback"
    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "os"
    "path/filepath"
)

// InitCommandOptions holds the required data for the `config init` command.
type InitCommandOptions struct {
    Overwrite bool
    DestDir   string
}

const defaultFileName = "config.toml"

func newConfigInitCommand() *cobra.Command {
    o := &InitCommandOptions{}

    initCmd := &cobra.Command{
        Use:   "init",
        Short: "Writes current configuration to a configuration file.",
        Long:  `Creates or updates the configuration file in the data directory or custom directory with the current configuration settings.`,
        Args:  cobra.NoArgs,
        Run:   o.Run,
    }

    initCmd.Flags().BoolVar(&o.Overwrite, "overwrite", false, "Overwrite existing configuration file.")
    initCmd.Flags().StringVar(&o.DestDir, "dest-dir", "", "Sets where to save the configuration file.")

    return initCmd
}

func (o *InitCommandOptions) Run(cobraCmd *cobra.Command, _ []string) {
    logrus.Info("Executing `gphotos-cli config init`")

    if o.DestDir == "" {
        o.DestDir = configuration.Settings.GetString("directories.Data")
    }

    absPath, err := filepath.Abs(o.DestDir)
    if err != nil {
        errMsg := fmt.Sprintf("Cannot find absolute path: %v", err)
        feedback.Fatal(errMsg, feedback.ErrGeneric)
    }
    configFileAbsPath := filepath.Join(absPath, defaultFileName)

    if !o.Overwrite && checkFileExists(configFileAbsPath) {
        feedback.Fatal("Config file already exists, use --overwrite to discard the existing one.", feedback.ErrGeneric)
    }

    logrus.Infof("Writing config file to: %s", absPath)

    if err := os.MkdirAll(absPath, os.FileMode(0755)); err != nil {
        errMsg := fmt.Sprintf("Cannot create config file directory: %v", err)
        feedback.Fatal(errMsg, feedback.ErrGeneric)
    }

    newSettings := viper.New()
    configuration.SetDefaults(newSettings)
    configuration.BindFlags(cobraCmd, newSettings)

    ff := configuration.FolderUpload{
        Path:              "PATH_TO_YOUR_PICTURES",
        CreateAlbums:      "folderName",
        DeleteAfterUpload: false,
        Include:           []string{},
        Exclude:           []string{},
    }

    ffa := []configuration.FolderUpload{ff}

    // Add a sample of the [[folders]] section.
    newSettings.SetDefault("folders", ffa)
	
    if err := newSettings.WriteConfigAs(configFileAbsPath); err != nil {
        errMsg := fmt.Sprintf("Cannot create config file: %v", err)
        feedback.Fatal(errMsg, feedback.ErrGeneric)
    }

    msg := fmt.Sprintf("Config file written to: %s", configFileAbsPath)
    logrus.Info(msg)
    feedback.Print(msg)
}

func checkFileExists(filePath string) bool {
    _, err := os.Stat(filePath)
    return err == nil
}
