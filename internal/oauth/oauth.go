package oauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/int128/oauth2cli"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/sync/errgroup"
)

var (
	// GoogleAuthEndpoint is the Google authentication endpoint.
	GoogleAuthEndpoint = google.Endpoint

	// PhotosLibraryScope is Google Photos OAuth2 scope.
	PhotosLibraryScope = "https://www.googleapis.com/auth/photoslibrary"

	ErrTokenIsNil = errors.New("OAuth 2.0 token is nil")
)

// Config represents a config for the OAuth 2.0 flow.
type Config struct {
	// OAuth's application ID.
	ClientID string
	// OAuth's application secret.
	ClientSecret string

	// Logger function for debug.
	Logf func(format string, args ...interface{})

	oAuth2Config *oauth2.Config
}

// GetToken refresh the provided token or create a new OAuth 2.0 token if the provided one is nil.
func GetToken(ctx context.Context, config *Config) (*oauth2.Token, error) {
	if err := config.validateAndSetDefaults(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config.getTokenFromWeb(ctx)
}

// RefreshToken refresh the provided token if needed.
func RefreshToken(ctx context.Context, config *Config, token *oauth2.Token) (*oauth2.Token, error) {
	if err := config.validateAndSetDefaults(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return config.refreshToken(ctx, token)
}

// Client returns an authenticated client using the specified token.
func Client(ctx context.Context, config *Config, token *oauth2.Token) (*http.Client, error) {
	if err := config.validateAndSetDefaults(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return config.oAuth2Config.Client(ctx, token), nil
}

func (c *Config) validateAndSetDefaults() error {
	if c.ClientID == "" || c.ClientSecret == "" {
		return fmt.Errorf("both ClientID and ClientSecret must be set")
	}

	if c.Logf == nil {
		c.Logf = func(string, ...interface{}) {}
	}

	c.oAuth2Config = &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Scopes:       []string{PhotosLibraryScope},
		Endpoint:     GoogleAuthEndpoint,
	}

	return nil
}

// getTokenFromWeb starts a local HTTP server, opens the web browser to initiate the OAuth Web application
// flow, blocks until the user completes authorization and is redirected back, and returns the access token.
func (c *Config) getTokenFromWeb(ctx context.Context) (*oauth2.Token, error) {
	ready := make(chan string, 1)
	cfg := oauth2cli.Config{
		OAuth2Config:         *c.oAuth2Config,
		LocalServerReadyChan: ready,
		Logf:                 c.Logf,
	}

	var token *oauth2.Token
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case url := <-ready:
			fmt.Printf("\nVisit the URL below in a browser:\n\n%s\n\n", url)
			fmt.Printf("If you are opening the url manually on a different machine you will need to curl the result URL on this machine manually.\n\n")
			return nil
		case <-ctx.Done():
			return fmt.Errorf("context done while waiting for authorization: %w", ctx.Err())
		}
	})
	eg.Go(func() error {
		tk, err := oauth2cli.GetToken(ctx, cfg)
		if err != nil {
			return fmt.Errorf("unable to get token: %w", err)
		}
		token = tk
		return nil
	})

	return token, eg.Wait()
}

// refreshToken refresh the OAuth 2.0 token.
func (c *Config) refreshToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	if token == nil {
		return nil, ErrTokenIsNil
	}

	if token.Valid() {
		return token, nil
	}

	return c.oAuth2Config.TokenSource(ctx, token).Token()
}
