package tokenstore

import (
	"encoding/json"

	"github.com/palantir/stacktrace"
	"golang.org/x/oauth2"

	keyring "github.com/zalando/go-keyring"
)

const (
	serviceName = "googlephotos-uploader-go-api"
)

var (
	ErrNotFound     = stacktrace.Propagate(keyring.ErrNotFound, "failed retrieving token from keyring")
	ErrInvalidToken = stacktrace.NewError("invalid token")
)

// StoreToken lets you store a token in the OS keyring
func StoreToken(googleUserEmail string, token *oauth2.Token) error {
	tokenJSONBytes, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = keyring.Set(serviceName, googleUserEmail, string(tokenJSONBytes))
	if err != nil {
		return stacktrace.Propagate(err, "failed storing token into keyring")
	}
	return nil
}

// RetrieveToken lets you get a token by google account email
func RetrieveToken(googleUserEmail string) (*oauth2.Token, error) {
	tokenJSONString, err := keyring.Get(serviceName, googleUserEmail)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed retrieving token from keyring")
	}

	var token oauth2.Token
	err = json.Unmarshal([]byte(tokenJSONString), &token)
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed unmarshaling token")
	}

	// validate token
	{
		if !token.Valid() {
			return nil, ErrInvalidToken
		}
	}

	return &token, nil
}
