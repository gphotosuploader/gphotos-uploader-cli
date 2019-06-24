package tokenstore

import (
	"encoding/json"
	"github.com/99designs/keyring"
	"golang.org/x/oauth2"
)

// KeyringRepository represents a repository provided by different secrets
// backend using github.com/99designs/keyring package.
//
// See https://github.com/99designs/keyring for details.
type KeyringRepository struct {
	keyring.Keyring
}

// NewKeyringRepository creates a new repository
func NewKeyringRepository() (*KeyringRepository, error) {
	kr, err := keyring.Open(keyring.Config{
		ServiceName: serviceName,
	})
	if err != nil {
		return nil, err
	}
	return &KeyringRepository{kr}, nil
}

// StoreToken lets you store a token in the OS keyring
func (r *KeyringRepository) StoreToken(email string, token *oauth2.Token) error {
	tokenJSONBytes, err := json.Marshal(token)
	if err != nil {
		return ErrInvalidToken
	}

	err = r.Set(keyring.Item{Key: email, Data: tokenJSONBytes})
	if err != nil {
		return ErrNotStoredToken
	}
	return nil
}

// RetrieveToken lets you get a token from the OS keyring.
// If the Token is not valid returns a ErrInvalidToken.
func (r *KeyringRepository) RetrieveToken(email string) (*oauth2.Token, error) {
	tk, err := r.getToken(email)
	if err != nil {
		return nil, err
	}

	// validate token
	if !tk.Valid() {
		return nil, ErrInvalidToken
	}

	return &tk, nil
}

// getToken returns the specified token from the repository
func (r *KeyringRepository) getToken(email string) (oauth2.Token, error) {
	item, err := r.Get(email)
	if err != nil {
		return oauth2.Token{}, ErrNotFound
	}

	var tk oauth2.Token
	err = json.Unmarshal(item.Data, &tk)
	if err != nil {
		return oauth2.Token{}, ErrInvalidToken
	}

	return tk, nil
}
