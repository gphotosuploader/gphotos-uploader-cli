package auth

import (
	"context"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/configuration"

	"github.com/spf13/cobra"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
)

// AuthCmd holds the required data for the init cmd
type AuthCmd struct {
}

func NewCommand() *cobra.Command {
	cmd := &AuthCmd{}

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
	cli, err := app.StartServices(ctx, configuration.Settings.GetString("directories.data"))
	if err != nil {
		return err
	}
	defer func() {
		_ = cli.Stop()
	}()

	_, err = cli.AuthenticateFromWeb(ctx)
	if err == nil {
		cli.Logger.Donef("Successful authentication for account '%s'", configuration.Settings.GetString("auth.account"))
	}

	return err
}
