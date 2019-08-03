package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/int128/oauth2cli"
	"golang.org/x/oauth2"

	"github.com/nmrshll/gphotos-uploader-cli/datastore/tokenstore"
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
func newOAuth2Client(ctx context.Context, tkm *tokenstore.Service, oauth2Config oauth2.Config, account string) (*http.Client, error) {
	token, err := tkm.RetrieveToken(account)
	if err != nil {
		log.Printf("Token has not been retrieved from token store: %s", err)
	}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, newHTTPClient())
	switch {
	case token == nil:
		token, err = oauth2cli.GetToken(ctx, oauth2cli.Config{
			OAuth2Config: oauth2Config,
			ShowLocalServerURL: func(url string) {
				log.Printf("Open %s", url)
			},
			LocalServerSuccessHTML: successPage,
		})
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
	err = tkm.StoreToken(account, token)
	if err != nil {
		return nil, fmt.Errorf("failed storing token: %s", err)
	}

	client := oauth2Config.Client(ctx, token)
	return client, nil
}
