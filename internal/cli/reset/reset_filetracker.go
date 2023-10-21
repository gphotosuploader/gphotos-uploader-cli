package reset

import (
	"context"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/feedback"
	"github.com/spf13/cobra"
)

// FileTrackerCommandOptions contains the input to the 'reset file-tracker' command.
type FileTrackerCommandOptions struct {
	*flags.GlobalFlags

	Force bool
}

func initFileTrackerCommand(globalFlags *flags.GlobalFlags) *cobra.Command {
	o := &FileTrackerCommandOptions{
		GlobalFlags: globalFlags,

		Force: false,
	}

	command := &cobra.Command{
		Use:   "file-tracker",
		Short: "Reset the already uploaded files database",
		Long:  `Reset the internal database which keep track of the already uploaded files.`,
		Args:  cobra.NoArgs,
		RunE:  o.Run,
	}

	command.Flags().BoolVar(&o.Force, "force", false, "Force the deletion without asking.")

	return command
}

func (o *FileTrackerCommandOptions) Run(cobraCmd *cobra.Command, args []string) error {
	ctx := context.Background()
	cli, err := app.Start(ctx, o.CfgDir)
	if err != nil {
		return err
	}
	defer func() {
		_ = cli.Stop()
	}()

	cli.Logger.Debug("Removing the File Tracker database...")

	// If the force flag is not user, ask for user confirmation
	if !o.Force {
		userConfirmation, err := feedback.YesNoPrompt("Do you want to reset the already uploaded file tracker?", false)
		if err != nil {
			return err
		}

		if !userConfirmation {
			cobraCmd.Println("User aborted the removal of the File Tracker")
			return nil
		}
	}

	err = cli.FileTracker.Destroy()
	if err != nil {
		return err
	}

	cobraCmd.Println("File Tracker was reset")

	return nil
}
