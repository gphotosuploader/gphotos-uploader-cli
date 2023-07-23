package config

import (
    "fmt"
    "github.com/gphotosuploader/gphotos-uploader-cli/internal/configuration"
    "github.com/gphotosuploader/gphotos-uploader-cli/internal/feedback"
    "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "gopkg.in/yaml.v3"
)

func initDumpCommand() *cobra.Command {
    var dumpCommand = &cobra.Command{
        Use:   "dump",
        Short: "Prints the current configuration",
        Long:  "Prints the current configuration.",
        Args:  cobra.NoArgs,
        Run:   runDumpCommand,
    }
    return dumpCommand
}

func runDumpCommand(_ *cobra.Command, _ []string) {
    logrus.Info("Executing `gphotos-cli config dump`")
    feedback.PrintResult(dumpResult{configuration.Settings.AllSettings()})
}

// output from this command requires special formatting, let's create a dedicated
// feedback.Result implementation
type dumpResult struct {
    data map[string]interface{}
}

func (dr dumpResult) Data() interface{} {
    // TODO: Hide sensible information: client_secret.
    
    return dr.data
}

func (dr dumpResult) String() string {
    bs, err := yaml.Marshal(dr.data)
    if err != nil {
        // Should never happen
        errMsg := fmt.Sprintf("unable to marshal config to YAML: %v", err)
        panic(errMsg)
    }
    return string(bs)
}
