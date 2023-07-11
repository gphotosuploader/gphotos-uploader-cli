package auth

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/flags"
)

// AuthCmd holds the required data for the init cmd
type AuthCmd struct {
	*flags.GlobalFlags
}

func NewCommand(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &AuthCmd{GlobalFlags: globalFlags}

	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate account with Google Photos to get OAuth 2.0 token",
		Long:  `Force authentication against Google Photos to get OAuth 2.0 token.`,
		Args:  cobra.NoArgs,
		RunE:  cmd.Run,
	}

	return authCmd
}

func (cmd *AuthCmd) Run(cobraCmd *cobra.Command, args []string) error {
	ctx := context.Background()
	cli, err := app.StartServices(ctx, cmd.CfgDir)
	if err != nil {
		return err
	}
	defer func() {
		_ = cli.Stop()
	}()

	_, err = cli.AuthenticateFromWeb(ctx)
	if err == nil {
		cli.Logger.Donef("Successful authentication for account '%s'", cli.Config.Account)
	}

	return err
}
