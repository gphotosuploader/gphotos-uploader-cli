package tokenstore

import (
	"errors"
)

var (
	// ErrNotFound is the expected error if the token is not found in the keyring
	ErrNotFound = errors.New("failed retrieving token from keyring")

	// ErrInvalidToken is the expected error if the token is not a valid one
	ErrInvalidToken = errors.New("invalid token")
)

