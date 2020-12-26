package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cmd/flags"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/photos"
)

// InitCmd holds the required data for the init cmd
type AuthCmd struct {
	*flags.GlobalFlags
}

func NewAuthCmd(globalFlags *flags.GlobalFlags) *cobra.Command {
	cmd := &AuthCmd{GlobalFlags: globalFlags}

	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with Google Photos to refresh tokens",
		Long:  `Force authentication against Google Photos to refresh tokens.`,
		Args:  cobra.NoArgs,
		RunE:  cmd.Run,
	}

	return authCmd
}

func (cmd *AuthCmd) Run(cobraCmd *cobra.Command, args []string) error {
	cli, err := app.Start(cmd.CfgDir)
	if err != nil {
		return err
	}
	defer func() {
		_ = cli.Stop()
	}()

	// get OAuth2 Configuration with our App credentials
	oauth2Config := oauth2.Config{
		ClientID:     cli.Config.APIAppCredentials.ClientID,
		ClientSecret: cli.Config.APIAppCredentials.ClientSecret,
		Endpoint:     photos.Endpoint,
		Scopes:       photos.Scopes,
	}

	ctx := context.Background()

	if _, err := cli.NewOAuth2Client(ctx, oauth2Config, cli.Config.Account); err != nil {
		cli.Logger.Failf("Failed authentication for account: %s", cli.Config.Account)
		cli.Logger.Debugf("Authentication error: err=%s", err)
		return err
	}
	cli.Logger.Donef("Successful authentication for account: %s", cli.Config.Account)

	return nil
}
