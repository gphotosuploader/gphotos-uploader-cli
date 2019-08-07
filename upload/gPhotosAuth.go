package upload

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	gphotosns "github.com/gphotosuploader/google-photos-api-client-go/noserver-gphotos"
	"github.com/juju/errors"
	"golang.org/x/oauth2"

	"github.com/gphotosuploader/gphotos-uploader-cli/datastore/tokenstore"
)

// Authenticate returns a Google Photos client form the supplied Google account.
// It will try to get the credentials from the Token Manager, if they are not valid will open an URL to get it
// from Google.
func Authenticate(tkm *tokenstore.Service, oauthConfig *oauth2.Config, account string) (*gphotosns.Client, error) {
	// try to load token from keyring
	token, err := tkm.RetrieveToken(account)
	if err == nil && token != nil { // if error ignore and skip
		// if found create client from token
		client, err := gphotosns.NewClient(gphotosns.FromToken(oauthConfig, token))
		if err == nil && client != nil { // if error ignore and skip
			return client, nil
		}
	}

	// else authenticate again to grab a new token
	fmt.Println(color.CyanString(fmt.Sprintf("Need to log login into account %s", account)))
	time.Sleep(1200 * time.Millisecond)
	client, err := gphotosns.NewClient(
		gphotosns.AuthenticateUser(
			oauthConfig,
			gphotosns.WithUserLoginHint(account),
		),
	)
	if err != nil {
		return nil, errors.Annotate(err, "failed authenticating new client")
	}

	// and store the token into the keyring
	err = tkm.StoreToken(account, client.Token())
	if err != nil {
		return nil, errors.Annotate(err, "failed storing token")
	}

	return client, nil
}
