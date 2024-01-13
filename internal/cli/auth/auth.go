package auth

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/cli/flags"
)

// AuthCmd holds the required data for the init cmd
type AuthCmd struct {
	*flags.GlobalFlags

	// command flags
	Port                int
	LocalBindAddress    string
	RedirectURLHostname string
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

	authCmd.Flags().IntVar(&cmd.Port, "port", 0, "port on which the auth server will listen (default 0)")
	authCmd.Flags().StringVar(&cmd.LocalBindAddress, "local-bind-address", "127.0.0.1", "local address on which the auth server will listen")
	authCmd.Flags().StringVar(&cmd.RedirectURLHostname, "redirect-url-hostname", "localhost", "hostname of the redirect URL")

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

	// customize authentication options based on the command line parameters
	authOptions := app.AuthenticationOptions{}
	authOptions.LocalServerBindAddress = fmt.Sprintf("%s:%d", cmd.LocalBindAddress, cmd.Port)

	if cmd.RedirectURLHostname != "" {
		authOptions.RedirectURLHostname = cmd.RedirectURLHostname
	}

	_, err = cli.AuthenticateFromWeb(ctx, authOptions)
	if err == nil {
		cli.Logger.Donef("Successful authentication for account '%s'", cli.Config.Account)
	}

	return err
}
