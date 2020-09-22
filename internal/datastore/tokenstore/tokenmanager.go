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
)

// TokenManager allows you to set/get oauth.Token into a permanent repository.
type TokenManager interface {
	// StoreToken stores a oauth.Token in the repository.
	StoreToken(email string, token *oauth2.Token) error

	// RetrieveToken returns the stored oauth.Token.
	RetrieveToken(email string) (*oauth2.Token, error)

	// Close closes the repository.
	Close() error
}
