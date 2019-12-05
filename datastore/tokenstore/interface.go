package tokenstore

import (
	"errors"

	"golang.org/x/oauth2"
)

var (
	// ErrNotFound is the expected error if the token is not found in the keyring
	ErrNotFound = errors.New("failed retrieving token from keyring")

	// ErrInvalidToken is the expected error if the token is not a valid one
	ErrInvalidToken = errors.New("invalid token")

	// ErrNotStoredToken is the expected error if the token has not been stored
	ErrNotStoredToken = errors.New("failed storing token into keyring")
)

const (
	serviceName = "gPhotosUploader"
)

type TokenManager interface {
	StoreToken(email string, token *oauth2.Token) error
	RetrieveToken(email string) (*oauth2.Token, error)
	Close() error
}
