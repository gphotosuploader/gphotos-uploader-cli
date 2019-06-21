package tokenstore

import (
	"encoding/json"
	"fmt"
	"github.com/99designs/keyring"
	"golang.org/x/oauth2"
	"log"
)

const (
	serviceName = "googlephotos-uploader-go-api"
)

var (
	// ErrNotFound is the expected error if the token isn't found in the keyring
	ErrNotFound = fmt.Errorf("failed retrieving token from keyring")

	// ErrInvalidToken is the expected error if the token isn't a valid one
	ErrInvalidToken = fmt.Errorf("invalid token")
)

var kr keyring.Keyring

func init() {
	// Use the best keyring implementation for your operating system
	k, err := keyring.Open(keyring.Config{
		ServiceName: serviceName,
	})
	if err != nil {
		log.Fatalln(err)
	}
	kr = k
}

// StoreToken lets you store a token in the OS keyring
func StoreToken(googleUserEmail string, token *oauth2.Token) error {
	tokenJSONBytes, err := json.Marshal(token)
	if err != nil {
		return err
	}

	//err = keyring.Set(serviceName, googleUserEmail, string(tokenJSONBytes))
	err = kr.Set(keyring.Item{Key: googleUserEmail, Data: tokenJSONBytes})
	if err != nil {
		return fmt.Errorf("failed storing token into keyring: %v", err)
	}
	return nil
}

// RetrieveToken lets you get a token by google account email
func RetrieveToken(googleUserEmail string) (*oauth2.Token, error) {
	//tokenJSONString, err := keyring.Get(serviceName, googleUserEmail)
	/*if err != nil {
		if err == keyring.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}*/
	item, err := kr.Get(googleUserEmail)
	if err != nil {
		log.Fatalln(err)
	}
	tokenJSONString := string(item.Data)

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
//func MockInit() {
//	keyring.MockInit()
//}
