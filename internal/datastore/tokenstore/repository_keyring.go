package tokenstore

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/99designs/keyring"
	"golang.org/x/crypto/ssh/terminal"
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
// backend could be used to select which backed will be used. If it's empty
// the library will select the most suitable depending OS.
//
// All currently supported secure storage backends:
//
//	SecretServiceBackend BackendType = "secret-service"
//	KeychainBackend      BackendType = "keychain"
//	KWalletBackend       BackendType = "kwallet"
//	WinCredBackend       BackendType = "wincred"
//	FileBackend          BackendType = "file"
//	PassBackend          BackendType = "pass"
func NewKeyringRepository(backend string, promptFunc *keyring.PromptFunc, keyringDir string) (*KeyringRepository, error) {
	keyringConfig := defaultConfig(keyringDir)
	if backend != "" && backend != "auto" {
		keyringConfig.AllowedBackends = append(keyringConfig.AllowedBackends, keyring.BackendType(backend))
	}
	if promptFunc != nil {
		keyringConfig.FilePasswordFunc = *promptFunc
	}
	kr, err := keyring.Open(keyringConfig)
	if err != nil {
		return nil, err
	}
	return &KeyringRepository{kr}, nil
}

func defaultConfig(keyringDir string) keyring.Config {
	return keyring.Config{
		ServiceName:          serviceName,
		KeychainName:         serviceName,
		KeychainPasswordFunc: encryptionKeyFromEnvVarOrTerminal,
		FilePasswordFunc:     encryptionKeyFromEnvVarOrTerminal,
		FileDir:              keyringDir,
	}
}

// StoreToken lets you store a token in the OS keyring
func (r *KeyringRepository) StoreToken(email string, token *oauth2.Token) error {
	if token.RefreshToken == "" {
		// Restore refresh token from previously stored token if available
		oldToken, _ := r.RetrieveToken(email)
		if oldToken != nil {
			token.RefreshToken = oldToken.RefreshToken
		}
	}

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
func (r *KeyringRepository) RetrieveToken(email string) (*oauth2.Token, error) {
	tk, err := r.getToken(email)
	if err != nil {
		return nil, err
	}

	// validate token.
	// We are not using token.Valid() because we want to retrieve token event though it's expired.
	if tk.AccessToken == "" {
		return nil, ErrInvalidToken
	}

	return &tk, nil
}

// Close closes the keyring repository.
func (r *KeyringRepository) Close() error {
	// in this particular implementation we don't need to do anything.
	return nil
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

// encryptionKeyFromEnvVarOrTerminal returns the key to open the keyring.
// It will read it from an environment var if is set, or read from the terminal otherwise.
func encryptionKeyFromEnvVarOrTerminal(_ string) (string, error) {
	if key := os.Getenv("GPHOTOS_CLI_TOKENSTORE_KEY"); len(key) > 0 {
		return key, nil
	}
	fmt.Print("Enter the passphrase to open the token store: ")
	b, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println()
	return string(b), nil
}
