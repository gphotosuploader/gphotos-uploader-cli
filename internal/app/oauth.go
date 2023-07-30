package app

import (
	"context"
	"fmt"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/configuration"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/oauth"
)

// AuthenticateFromToken returns an HTTP client authenticated in Google Photos.
// AuthenticateFromToken will use the token from the Token Manage.
func (app *App) AuthenticateFromToken(ctx context.Context) (*http.Client, error) {
	account := configuration.Settings.GetString("auth.account")
	logrus.Infof("Authenticating using token for '%s'", account)

	token, err := app.TokenManager.Get(account)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token, have you authenticated before?: %w", err)
	}

	cfg := &oauth.Config{
		ClientID:     configuration.Settings.GetString("auth.client_id"),
		ClientSecret: configuration.Settings.GetString("auth.client_secret"),
		Logf:         app.Logger.Debugf,
	}

	token, err = oauth.RefreshToken(ctx, cfg, token)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Token is valid, expires at %s", token.Expiry)

	if err := app.TokenManager.Put(account, token); err != nil {
		logrus.Debugf("Failed to store token into token manager: %s", err)
	}

	return oauth.Client(ctx, cfg, token)
}

// AuthenticateFromWeb returns an HTTP client authenticated in Google Photos.
// AuthenticateFromWeb will create a new token after completing the OAuth 2.0 flow.
func (app *App) AuthenticateFromWeb(ctx context.Context) (*http.Client, error) {
	account := configuration.Settings.GetString("auth.account")
	logrus.Infof("Getting authentication token for '%s'", account)

	cfg := &oauth.Config{
		ClientID:     configuration.Settings.GetString("auth.client_id"),
		ClientSecret: configuration.Settings.GetString("auth.client_secret"),
		Logf:         logrus.Debugf,
	}

	token, err := oauth.GetToken(ctx, cfg)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Token obtained, expires at %s", token.Expiry)

	if err := app.TokenManager.Put(account, token); err != nil {
		logrus.Debugf("Failed to store token into token manager: %s", err)
	}

	return oauth.Client(ctx, cfg, token)
}
