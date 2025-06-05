package auth

import (
	"context"
	"fmt"
	"net/url"

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
	RedirectURL         string
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

	// Declare the version flag and then you can deprecate it.
	authCmd.Flags().StringVar(&cmd.RedirectURLHostname, "redirect-url-hostname", "", "hostname of the redirect URL")
	err := authCmd.Flags().MarkDeprecated("redirect-url-hostname", "use --redirect-url instead")
	if err != nil {
		panic(fmt.Sprintf("error marking flag --redirect-url-hostname as deprecated: %v", err))
	}
	authCmd.Flags().StringVar(&cmd.RedirectURL, "redirect-url", "", "URL of the redirect URL to use for authentication, (e.g. http://localhost:12345)")

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

	// TODO: Remove this check in the v6 release.
	// --redirect-url-hostname is deprecated, so we keep it for backward compatibility.
	if cmd.RedirectURLHostname != "" && cmd.Port == 0 {
		return fmt.Errorf("--port is required when using --redirect-url-hostname")
	}
	if cmd.RedirectURL != "" && cmd.RedirectURLHostname != "" {
		return fmt.Errorf("--redirect-url and --redirect-url-hostname cannot be used together")
	}
	// End of TODO

	// Validate the redirect URL if it is set.
	if cmd.RedirectURL != "" {
		_, err := url.ParseRequestURI(cmd.RedirectURL)
		if err != nil {
			return fmt.Errorf("invalid redirect URL '%s'", cmd.RedirectURL)
		}
	}

	// If redirect URL is set, we require the port to be specified as well.
	if cmd.RedirectURL != "" && cmd.Port == 0 {
		return fmt.Errorf("--port is required when using --redirect-url")
	}

	// customize authentication options based on the command line parameters
	authOptions := app.AuthenticationOptions{}
	authOptions.LocalServerBindAddress = fmt.Sprintf("%s:%d", cmd.LocalBindAddress, cmd.Port)

	// TODO: Remove this check in the v6 release.
	// If redirect URL hostname is set, we use it to construct the redirect URL.
	if cmd.RedirectURLHostname != "" {
		authOptions.RedirectURL = fmt.Sprintf("http://%s:%d", cmd.RedirectURLHostname, cmd.Port)
	}
	// End of TODO

	// If redirect URL is set, we use it to construct the redirect URL.
	if cmd.RedirectURL != "" {
		authOptions.RedirectURL = cmd.RedirectURL
	}

	_, err = cli.AuthenticateFromWeb(ctx, authOptions)
	if err == nil {
		cli.Logger.Donef("Successful authentication for account '%s'", cli.Config.Account)
	}

	return err
}
