package gphotosapiclient

import (
	"context"
	"net/http"

	oauth2ns "github.com/nmrshll/oauth2-noserver"
	"gitlab.com/nmrshll/gphotos-uploader-go-api/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	photoslibrary "google.golang.org/api/photoslibrary/v1"
)

var (
	OAuthConfig = oauth2.Config{
		ClientID:     config.API_APP_CREDENTIALS.ClientID,
		ClientSecret: config.API_APP_CREDENTIALS.ClientSecret,
		Scopes:       []string{photoslibrary.PhotoslibraryScope},
		Endpoint:     google.Endpoint,
	}
)

// NewOAuthClient creates a new http.Client with a bearer access token
// TODO: refactor to load apiAppCredentials from config file
func NewOAuthClient() (*oauth2ns.AuthorizedClient, error) {
	// conf := &oauth2.Config{
	// 	ClientID:     clientID,
	// 	ClientSecret: clientSecret,
	// 	Scopes:       []string{photoslibrary.PhotoslibraryScope},
	// 	Endpoint:     google.Endpoint,
	// }
	photosClient := oauth2ns.Authorize(&OAuthConfig)

	return photosClient, nil
}

func NewClientFromToken(token *oauth2.Token) *http.Client {
	return OAuthConfig.Client(context.Background(), token)
}
