package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/int128/oauth2cli"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"

	"github.com/gphotosuploader/gphotos-uploader-cli/app"
)

const successPage = `
		<div style="height:100px; width:100%!; display:flex; flex-direction: column; justify-content: center; align-items:center; background-color:#2ecc71; color:white; font-size:22"><div>Success!</div></div>
		<p style="margin-top:20px; font-size:18; text-align:center">You are authenticated, you can now return to the program. This will auto-close</p>
		<script>window.onload=function(){setTimeout(this.close, 4000)}</script>
		`

func newHTTPClient() *http.Client {
	return http.DefaultClient
}

// newOAuth2Client returns a http client for the supplied Google account.
// It will try to get the credentials from the Token Manager, if they are not valid will try to refresh the token or
// ask for authenticate again.
func newOAuth2Client(ctx context.Context, tokenManager app.TokenManager, oauth2Config oauth2.Config, account string) (*http.Client, error) {
	token, err := tokenManager.RetrieveToken(account)
	if err != nil {
		log.Printf("Token has not been retrieved from token store: %s", err)
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, newHTTPClient())
	switch {
	case token == nil:
		token, err = obtainOAuthTokenFromAuthServer(ctx, oauth2Config)
		if err != nil {
			return nil, fmt.Errorf("could not get a token: %s", err)
		}

	case !token.Valid():
		log.Printf("Token has been expired, refreshing")
		token, err = oauth2Config.TokenSource(ctx, token).Token()
		if err != nil {
			return nil, fmt.Errorf("could not refresh the token: %s", err)
		}
	}

	// debug
	if token != nil {
		log.Printf("Token expiration: %s", token.Expiry.String())
	}

	// and store the token into the keyring
	err = tokenManager.StoreToken(account, token)
	if err != nil {
		return nil, fmt.Errorf("failed storing token: %s", err)
	}

	client := oauth2Config.Client(ctx, token)
	return client, nil
}

func obtainOAuthTokenFromAuthServer(ctx context.Context, oauth2Config oauth2.Config) (*oauth2.Token, error) {
	var token *oauth2.Token
	var err error

	ready := make(chan string, 1)
	var eg errgroup.Group
	eg.Go(func() error {
		select {
		case url, ok := <-ready:
			if !ok {
				return nil
			}
			// Open a browser to complete OAuth process.
			log.Printf("Openning browser to complete authorization.")
			err = browser.OpenURL(url)
			if err != nil {
				log.Printf("Browser was not detected. Complete the authorization browsing to: %s", url)
			}
			return nil
		case err := <-ctx.Done():
			return fmt.Errorf("context done while waiting for authorization: %s", err)
		}
	})
	eg.Go(func() error {
		defer close(ready)
		token, err = oauth2cli.GetToken(ctx, oauth2cli.Config{
			OAuth2Config:           oauth2Config,
			LocalServerReadyChan:   ready,
			LocalServerSuccessHTML: successPage,
		})
		return err
	})
	if err := eg.Wait(); err != nil {
		log.Printf("error while authorization: %s", err)
	}

	return token, err
}
