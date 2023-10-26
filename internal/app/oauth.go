package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/oauth"
)

// AuthenticateFromToken returns an HTTP client authenticated in Google Photos.
// AuthenticateFromToken will use the token from the Token Manage.
func (app *App) AuthenticateFromToken(ctx context.Context) (*http.Client, error) {
	account := app.Config.Account
	app.Logger.Infof("Authenticating using token for '%s'", account)

	token, err := app.TokenManager.Get(account)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token, have you authenticated before?: %w", err)
	}

	cfg := &oauth.Config{
		ClientID:     app.Config.APIAppCredentials.ClientID,
		ClientSecret: app.Config.APIAppCredentials.ClientSecret,
		Logf:         app.Logger.Debugf,
	}

	token, err = oauth.RefreshToken(ctx, cfg, token)
	if err != nil {
		return nil, err
	}

	app.Logger.Donef("Token is valid, expires at %s", token.Expiry)

	if err := app.TokenManager.Put(account, token); err != nil {
		app.Logger.Debugf("Failed to store token into token manager: %s", err)
	}

	return oauth.Client(ctx, cfg, token)
}

type AuthenticationOptions struct {
	// Hostname of the redirect URL.
	// You can set this if your provider does not accept localhost.
	// Default to localhost.
	RedirectURLHostname string

	// Hostname and port which the local server binds to.
	// You can set port number to 0 to allocate a free port.
	// If nil or an empty slice is given, it defaults to "127.0.0.1:0" i.e., a free port.
	LocalServerBindAddress string
}

// AuthenticateFromWeb returns an HTTP client authenticated in Google Photos.
// AuthenticateFromWeb will create a new token after completing the OAuth 2.0 flow.
func (app *App) AuthenticateFromWeb(ctx context.Context, authOptions AuthenticationOptions) (*http.Client, error) {
	account := app.Config.Account
	app.Logger.Infof("Getting authentication token for '%s'", account)

	cfg := &oauth.Config{
		ClientID:               app.Config.APIAppCredentials.ClientID,
		ClientSecret:           app.Config.APIAppCredentials.ClientSecret,
		Logf:                   app.Logger.Debugf,
		LocalServerBindAddress: []string{authOptions.LocalServerBindAddress},
		RedirectURLHostname:    authOptions.RedirectURLHostname,
	}

	token, err := oauth.GetToken(ctx, cfg)
	if err != nil {
		return nil, err
	}

	app.Logger.Donef("Token obtained, expires at %s", token.Expiry)

	if err := app.TokenManager.Put(account, token); err != nil {
		app.Logger.Debugf("Failed to store token into token manager: %s", err)
	}

	return oauth.Client(ctx, cfg, token)
}
