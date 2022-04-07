package app

import (
	"context"
	"net/http"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/oauth"
)

// NewOAuth2Client returns a HTTP client authenticated in Google Photos.
// NewOAuth2Client will get (from Token Manager) or create the token.
func (app App) NewOAuth2Client(ctx context.Context) (*http.Client, error) {
	account := app.Config.Account
	app.Logger.Infof("Getting OAuth token for '%s'", account)

	token, err := app.TokenManager.Get(account)
	if err != nil {
		app.Logger.Debugf("Unable to retrieve token from token manager: %s", err)
	}

	cfg := &oauth.Config{
		ClientID:     app.Config.APIAppCredentials.ClientID,
		ClientSecret: app.Config.APIAppCredentials.ClientSecret,
		Logf:         app.Logger.Debugf,
	}

	token, err = oauth.GetToken(ctx, cfg, token)
	if err != nil {
		return nil, err
	}

	app.Logger.Donef("Token is valid, expires at %s", token.Expiry)

	if err := app.TokenManager.Put(account, token); err != nil {
		app.Logger.Debugf("Failed to store token into token manager: %s", err)
	}

	return oauth.Client(ctx, cfg, token)
}
