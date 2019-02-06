package tokenstore

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"

	"github.com/godbus/dbus"
	keyring "github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

const (
	serviceName = "googlephotos-uploader-go-api"
)

type TokenStoreInterface interface {
	StoreToken(googleUserEmail string, token *oauth2.Token) error
	RetrieveToken(googleUserEmail string) (*oauth2.Token, error)
}

// the active TokenStore for this instance
var TokenStore TokenStoreInterface

func KeyRingSupported() bool {
	if runtime.GOOS == "linux" {
		// test dbus connection
		_, err := dbus.SessionBus()
		if err != nil {
			log.Print("No Dbus support")
			return false
		}
	}
	log.Print("Keyring is supported")
	return true
}

var (
	// ErrNotFound is the expected error if the token isn't found in the keyring
	ErrNotFound = fmt.Errorf("failed retrieving token from keyring")

	// ErrInvalidToken is the expected error if the token isn't a valid one
	ErrInvalidToken = fmt.Errorf("invalid token")
)

// TokenStoreKeyring Default token store that uses the os-specific keyring (via zalondo/go-keyring)
type TokenStoreKeyring struct{}

// StoreToken lets you store a token in the OS keyring
func (t TokenStoreKeyring) StoreToken(googleUserEmail string, token *oauth2.Token) error {
	tokenJSONBytes, err := json.Marshal(token)
	if err != nil {
		return err
	}

	err = keyring.Set(serviceName, googleUserEmail, string(tokenJSONBytes))
	if err != nil {
		return fmt.Errorf("failed storing token into keyring: %v", err)
	}
	return nil
}

// RetrieveToken lets you get a token by google account email
func (t TokenStoreKeyring) RetrieveToken(googleUserEmail string) (*oauth2.Token, error) {
	tokenJSONString, err := keyring.Get(serviceName, googleUserEmail)
	if err != nil {
		if err == keyring.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal([]byte(tokenJSONString), &token)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshaling token: %v", err)
	}

	// validate token
	if !token.Valid() {
		return nil, ErrInvalidToken
	}

	return &token, nil
}

// MockInit sets the provider to a mocked memory store, using keyring mock
func MockInit() {
	keyring.MockInit()
}
