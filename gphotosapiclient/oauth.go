package gphotosapiclient

import (
	"context"
	"net/http"
	"net/url"

	"github.com/palantir/stacktrace"

	"github.com/nmrshll/gphotos-uploader-cli/config"
	oauth2ns "github.com/nmrshll/oauth2-noserver"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

func OAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.Cfg.APIAppCredentials.ClientID,
		ClientSecret: config.Cfg.APIAppCredentials.ClientSecret,
		Scopes:       []string{photoslibrary.PhotoslibraryScope},
		Endpoint:     google.Endpoint,
	}
}

// NewOAuthClient creates a new http.Client with a bearer access token
func NewOAuthClient() (*oauth2ns.AuthorizedClient, error) {
	values := url.Values{}
	values.Set("login_hint", "admsommer21@gmail.com")
	oauthClient, err := oauth2ns.Authorize(OAuthConfig(), oauth2ns.WithAuthCallHTTPParams(values))
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed authorizing application and creating oauth client")
	}

	return oauthClient, nil
}

// NewClientFromToken returns an authorized google photos client from a user token
func NewClientFromToken(token *oauth2.Token) *http.Client {
	return OAuthConfig().Client(context.Background(), token)
}
