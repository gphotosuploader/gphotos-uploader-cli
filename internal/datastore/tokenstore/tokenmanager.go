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

// StoreToken stores the Token into the repository.
func (tm *TokenManager) StoreToken(email string, token *oauth2.Token) error {
	if token.AccessToken == "" {
		return ErrInvalidToken
	}

	// Restore refresh token from previously stored token if it's not available on the current one.
	token.RefreshToken = tm.getRefreshToken(email, token)

	return tm.repo.Set(email, token)
}

// RetrieveToken returns the Token from the repository.
func (tm *TokenManager) RetrieveToken(email string) (*oauth2.Token, error) {
	tk, err := tm.repo.Get(email)
	if err != nil {
		return nil, err
	}

	return tk, nil
}

// Close closes the Token repository.
func (tm *TokenManager) Close() error {
	return tm.repo.Close()
}

// getRefreshToken returns the most updated Refresh Token for the account (email).
func (tm *TokenManager) getRefreshToken(email string, token *oauth2.Token) string {
	// Get the previous stored Token.
	oldToken, err := tm.repo.Get(email)
	if token.RefreshToken == "" && err == nil {
		return oldToken.RefreshToken
	}
	return token.RefreshToken
}
