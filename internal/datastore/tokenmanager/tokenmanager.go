package tokenmanager

import (
	"errors"

	"golang.org/x/oauth2"
)

var (
	// ErrTokenNotFound is the expected error if the token is not found.
	ErrTokenNotFound = errors.New("token was not found")

	// ErrInvalidToken is the expected error if the token is not a valid.
	ErrInvalidToken = errors.New("invalid token")
)

// TokenManager allows to store and retrieve Token from a repository.
type TokenManager struct {
	repo Repository
}

// Repository is responsible to store tokens somewhere.
type Repository interface {
	Get(key string) (*oauth2.Token, error)
	Set(key string, token *oauth2.Token) error
	Close() error
}

// New returns a TokenManager using the specified repo.
func New(repo Repository) *TokenManager {
	return &TokenManager{repo: repo}
}

// Put stores the token, using email as key, into the repository.
func (tm *TokenManager) Put(email string, token *oauth2.Token) error {
	if token.AccessToken == "" {
		return ErrInvalidToken
	}

	if token.RefreshToken == "" {
		// Set RefreshToken from the token previously stored.
		if previousToken, err := tm.repo.Get(email); err == nil {
			token.RefreshToken = previousToken.RefreshToken
		}
	}

	return tm.repo.Set(email, token)
}

// Get returns the token, associated with email, from the repository.
func (tm *TokenManager) Get(email string) (*oauth2.Token, error) {
	tk, err := tm.repo.Get(email)
	if err != nil {
		return nil, err
	}

	return tk, nil
}

// Close closes the token repository.
func (tm *TokenManager) Close() error {
	return tm.repo.Close()
}
