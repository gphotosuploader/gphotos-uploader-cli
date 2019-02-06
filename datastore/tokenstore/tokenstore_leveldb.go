package tokenstore

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"golang.org/x/oauth2"
)

const keyPrefix = "credential"

// TokenLevelDB Stores credentials in main leveldb under key "credential_USERNAME
type TokenLevelDB struct {
	DB *leveldb.DB
}

func keyFor(user string) []byte {
	return []byte(fmt.Sprintf("%s_%s", keyPrefix, user))
}

// RetrieveToken return users token
func (t TokenLevelDB) RetrieveToken(user string) (*oauth2.Token, error) {
	tokenJSONString, err := t.DB.Get(keyFor(user), nil)
	if err == leveldb.ErrNotFound {
		log.Printf("Error finding credential")
		return nil, ErrNotFound
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

// StoreToken set users token
func (t TokenLevelDB) StoreToken(user string, token *oauth2.Token) error {
	tokenJSONBytes, err := json.Marshal(token)
	if err != nil {
		log.Printf("error marshalling token")
		return err
	}

	err = t.DB.Put(keyFor(user), tokenJSONBytes, nil)

	if err != nil {
		return err
	}
	log.Printf("stored token for user: %s", user)
	return nil
}
