package tokenmanager

import (
	"encoding/json"
	"fmt"
	"os"
	"syscall"

	"github.com/99designs/keyring"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"
)

// KeyringRepository represents a repository provided by different secrets
// backend using `99designs/keyring` package.
type KeyringRepository struct {
	store keyring.Keyring
}

// defaultConfig returns the default configuration from the keyring package.
func defaultConfig(keyringDir string) keyring.Config {
	const serviceName = "gPhotosUploader"
	return keyring.Config{
		ServiceName:          serviceName,
		KeychainName:         serviceName,
		KeychainPasswordFunc: promptFn(&stdInPasswordReader{}),
		FilePasswordFunc:     promptFn(&stdInPasswordReader{}),
		FileDir:              keyringDir,
	}
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
	return &KeyringRepository{store: kr}, nil
}

// Set stores a token into the OS keyring.
func (r *KeyringRepository) Set(email string, token *oauth2.Token) error {
	tokenJSONBytes, err := json.Marshal(token)
	if err != nil {
		return ErrInvalidToken
	}

	return r.store.Set(keyring.Item{
		Key:  email,
		Data: tokenJSONBytes,
	})
}

// getToken returns the specified token from the repository.
func (r *KeyringRepository) Get(email string) (*oauth2.Token, error) {
	var nullToken = &oauth2.Token{}

	item, err := r.store.Get(email)
	if err != nil {
		return nullToken, ErrTokenNotFound
	}

	var token oauth2.Token
	if err := json.Unmarshal(item.Data, &token); err != nil {
		return nullToken, ErrInvalidToken
	}

	return &token, nil
}

// Close closes the keyring repository.
func (r *KeyringRepository) Close() error {
	// in this particular implementation we don't need to do anything.
	return nil
}

// passwordReader represents a function to read a password.
type passwordReader interface {
	ReadPassword() (string, error)
}

// promptFn returns the key to open the keyring.
// It will read it from an environment var if is set, or read from the terminal otherwise.
func promptFn(pr passwordReader) func(string) (string, error) {
	return func(_ string) (string, error) {
		if key := os.Getenv("GPHOTOS_CLI_TOKENSTORE_KEY"); len(key) > 0 {
			return key, nil
		}
		fmt.Print("Enter the passphrase to open the token store: ")
		fmt.Println()
		return pr.ReadPassword()
	}
}

// stdInPasswordReader reads a password from the stdin.
type stdInPasswordReader struct{}

func (pr *stdInPasswordReader) ReadPassword() (string, error) {
	pwd, err := terminal.ReadPassword(syscall.Stdin)
	return string(pwd), err
}
