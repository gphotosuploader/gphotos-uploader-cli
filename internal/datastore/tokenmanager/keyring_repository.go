package tokenmanager

import (
	"encoding/json"
	"github.com/99designs/keyring"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/feedback"
	"golang.org/x/oauth2"
	"os"
)

// KeyringRepository represents a repository provided by different secrets
// backend using `99designs/keyring` package.
type KeyringRepository struct {
	store keyring.Keyring
}

// defaultConfig returns the default configuration from the keyring package.
func defaultConfig(keyringDir string) keyring.Config {
	const serviceName = "GooglePhotosCLI"
	return keyring.Config{
		ServiceName:          serviceName,
		KeychainName:         serviceName,
		KeychainPasswordFunc: getPassphraseFromEnvOrUserInputFn(),
		FilePasswordFunc:     getPassphraseFromEnvOrUserInputFn(),
		FileDir:              keyringDir,
		AllowedBackends:      supportedBackendTypes,
	}
}

// The `99designs/keyring` package supports several secure storage backends,
// but this CLI have implemented just some of them:
//
//	SecretServiceBackend BackendType = "secret-service"
//	KeychainBackend      BackendType = "keychain"
//	KWalletBackend       BackendType = "kwallet"
//	FileBackend          BackendType = "file"
var supportedBackendTypes = []keyring.BackendType{
	// MacOS
	keyring.KeychainBackend,
	// Linux
	keyring.SecretServiceBackend,
	keyring.KWalletBackend,
	// General
	keyring.FileBackend,
}

// NewKeyringRepository creates a new repository
// backend could be used to select which backed will be used. If it's empty or auto,
// the library will select the most suitable depending on OS.
func NewKeyringRepository(backend string, promptFunc *keyring.PromptFunc, keyringDir string) (*KeyringRepository, error) {
	keyringConfig := defaultConfig(keyringDir)
	if backend != "" && backend != "auto" {
		keyringConfig.AllowedBackends = []keyring.BackendType{
			keyring.BackendType(backend),
		}
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

// Get returns the specified token from the repository.
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
	// in this particular implementation, we don't need to do anything.
	return nil
}

// getPassphraseFromEnvOrUserInputFn returns the key to open the keyring.
// It will read it from an environment var if it's set, or read from the terminal otherwise.
func getPassphraseFromEnvOrUserInputFn() func(string) (string, error) {
	return func(_ string) (string, error) {
		// TODO: Use the configuration package to gather this env var.
		if key, ok := os.LookupEnv("GPHOTOS_CLI_TOKENSTORE_KEY"); ok {
			return key, nil
		}

		return feedback.InputUserField("Enter the passphrase to open the token store: ", true)
	}
}
