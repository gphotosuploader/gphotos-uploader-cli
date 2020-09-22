package tokenstore

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/99designs/keyring"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
)

const (
	serviceName = "gPhotosUploader"
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

// StoreToken lets you store a token in the OS keyring
func (r *KeyringRepository) StoreToken(email string, token *oauth2.Token) error {
	if token.AccessToken == "" {
		return ErrInvalidToken
	}

	// Restore refresh token from previously stored token if it's not available on the current one.
	token.RefreshToken = r.getRefreshToken(email, token)

	return r.setToken(email, token)
}

// RetrieveToken lets you get a token from the OS keyring.
func (r *KeyringRepository) RetrieveToken(email string) (*oauth2.Token, error) {
	tk, err := r.getToken(email)
	if err != nil {
		return nil, err
	}

	return tk, nil
}

// Close closes the keyring repository.
func (r *KeyringRepository) Close() error {
	// in this particular implementation we don't need to do anything.
	return nil
}

// getRefreshToken returns the most updated Refresh Token for the account (email).
func (r *KeyringRepository) getRefreshToken(email string, token *oauth2.Token) string {
	if token.RefreshToken != "" {
		return token.RefreshToken
	}

	// Returns the previous Refresh Token for the account (email).
	if token, err := r.getToken(email); err == nil {
		return token.RefreshToken
	}
	return ""
}

func defaultConfig(keyringDir string) keyring.Config {
	return keyring.Config{
		ServiceName:          serviceName,
		KeychainName:         serviceName,
		KeychainPasswordFunc: promptFn(StdInPasswordReader{}),
		FilePasswordFunc:     promptFn(StdInPasswordReader{}),
		FileDir:              keyringDir,
	}
}

// setToken stores the token into the repository.
func (r *KeyringRepository) setToken(email string, token *oauth2.Token) error {
	tokenJSONBytes, err := json.Marshal(token)
	if err != nil {
		return ErrInvalidToken
	}

	return r.Set(keyring.Item{
		Key:  email,
		Data: tokenJSONBytes,
	})
}

// getToken returns the specified token from the repository.
func (r *KeyringRepository) getToken(email string) (*oauth2.Token, error) {
	var nullToken = &oauth2.Token{}

	item, err := r.Get(email)
	if err != nil {
		return nullToken, ErrNotFound
	}

	var token oauth2.Token
	if err := json.Unmarshal(item.Data, &token); err != nil {
		return nullToken, ErrInvalidToken
	}

	return &token, nil
}

// PasswordReader represents a function to read a password.
type PasswordReader interface {
	ReadPassword() (string, error)
}

// promptFn returns the key to open the keyring.
// It will read it from an environment var if is set, or read from the terminal otherwise.
func promptFn(pr PasswordReader) func(string) (string, error) {
	return func(_ string) (string, error) {
		if key := os.Getenv("GPHOTOS_CLI_TOKENSTORE_KEY"); len(key) > 0 {
			return key, nil
		}
		fmt.Print("Enter the passphrase to open the token store: ")
		fmt.Println()
		return pr.ReadPassword()
	}
}

// StdInPasswordReader reads a password from the stdin.
type StdInPasswordReader struct{}

func (pr StdInPasswordReader) ReadPassword() (string, error) {
	pwd, err := terminal.ReadPassword(syscall.Stdin)
	return string(pwd), err
}
