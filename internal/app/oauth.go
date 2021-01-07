package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type PromptFunc func(code string) string

var (
	// GoogleAuthEndpoint is the Google authentication endpoint.
	GoogleAuthEndpoint = google.Endpoint

	// PhotosLibraryScope is Google Photos OAuth2 scope.
	PhotosLibraryScope = "https://www.googleapis.com/auth/photoslibrary"

	// AskForAuthCodeFn is the function used to get the Authorization code.
	// Useful for testing
	AskForAuthCodeFn = askForAuthCodeInTerminal
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

	oauth2Config := oauth2.Config{
		ClientID:     app.Config.APIAppCredentials.ClientID,
		ClientSecret: app.Config.APIAppCredentials.ClientSecret,
		Scopes:       []string{PhotosLibraryScope},
		Endpoint:     GoogleAuthEndpoint,
	}

	switch {
	case token == nil:
		app.Logger.Debug("Getting OAuth2 token from prompt...")
		token, err = getOfflineOAuth2Token(ctx, oauth2Config)
		if err != nil {
			return nil, fmt.Errorf("unable to get token: %s", err)
		}

	case !token.Valid():
		app.Logger.Debug("Token has been expired, refreshing it...")
		token, err = oauth2Config.TokenSource(ctx, token).Token()
		if err != nil {
			app.Logger.Errorf("Unable to refresh the token, err: %s", err)
			return nil, fmt.Errorf("unable to refresh the token: %s", err)
		}
	}

	app.Logger.Donef("Token is valid, expires at %s", token.Expiry.String())

	if err := app.TokenManager.Put(account, token); err != nil {
		app.Logger.Debugf("Failed to store token into token manager: %s", err)
	}

	client := oauth2Config.Client(ctx, token)
	return client, nil
}

func getOfflineOAuth2Token(ctx context.Context, oauth2Config oauth2.Config) (*oauth2.Token, error) {
	oauth2Config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"

	// Redirect user to consent page to ask for permission for the specified scopes.
	url := oauth2Config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	code, err := AskForAuthCodeFn(os.Stdin, url)
	if err != nil {
		return nil, err
	}

	// Use the custom HTTP client with a short timeout when requesting a token.
	ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Timeout: 2 * time.Second})
	return oauth2Config.Exchange(ctx, code)
}

func askForAuthCodeInTerminal(r io.Reader, url string) (string, error) {
	fmt.Printf("\nVisit the following URL in your browser:\n%v\n\n", url)

	var code string
	fmt.Print("After completing the authorization flow, enter the authorization code here: ")
	n, err := fmt.Fscanln(r, &code)
	if err != nil || n == 0 {
		return "", fmt.Errorf("unable to read authorization code: %s", err)
	}
	return code, nil
}
