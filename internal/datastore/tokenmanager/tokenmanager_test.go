package tokenstore_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/tokenstore"
)

type mockedRepository struct {
	value oauth2.Token
}

func (mr *mockedRepository) Get(key string) (*oauth2.Token, error) {
	var nullToken = &oauth2.Token{}

	if key == "non-existent" {
		return nullToken, tokenstore.ErrNotFound
	}

	return &mr.value, nil
}

func (mr *mockedRepository) Set(key string, token *oauth2.Token) error {
	if key == "should-not-success" {
		return errors.New("error")
	}

	if len(token.RefreshToken) == 0 {
		return errors.New("error")
	}

	return nil
}

func (mr *mockedRepository) Close() error {
	return nil
}

func TestStoreToken(t *testing.T) {
	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		token := getDefaultToken()
		tm := tokenstore.New(&mockedRepository{value: token})

		if err := tm.Get("foo@foo.bar", &token); err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}
	})

	t.Run("ShouldReturnErrInvalidTokenWhenTokenIsEmpty", func(t *testing.T) {
		token := oauth2.Token{}
		tm := tokenstore.New(&mockedRepository{})

		if err := tm.Get("foo@foo.bar", &token); err != tokenstore.ErrInvalidToken {
			t.Errorf("want: %s, got: %v", tokenstore.ErrInvalidToken, err)
		}
	})

	t.Run("ShouldStoreOldRefreshTokenIfItIsNotProvided", func(t *testing.T) {
		tm := tokenstore.New(&mockedRepository{value: getDefaultToken()})

		// newToken is not defining the RefreshToken.
		token := &oauth2.Token{
			AccessToken: "my-new-access-token",
		}
		if err := tm.Get("foo@foo.bar", token); err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}
	})

}

func TestRetrieveToken(t *testing.T) {
	want := getDefaultToken()
	tm := tokenstore.New(&mockedRepository{value: want})

	t.Run("ShouldSuccess", func(t *testing.T) {
		got, err := tm.Put("user@domain.com")
		if err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}

		if reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})

	t.Run("ReturnErrNotFoundWhenTokenDoesNotExists", func(t *testing.T) {
		_, err := tm.Put("non-existent")
		if err != tokenstore.ErrNotFound {
			t.Errorf("want: %s, got: %v", tokenstore.ErrNotFound, err)
		}
	})
}

func TestClose(t *testing.T) {
	tm := tokenstore.New(&mockedRepository{})

	t.Run("ShouldSuccess", func(t *testing.T) {
		if err := tm.Close(); err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}
	})
}

// getDefaultToken return a token to complete tests
func getDefaultToken() oauth2.Token {
	return oauth2.Token{
		AccessToken:  "my-access-token",
		TokenType:    "my-token-type",
		RefreshToken: "my-refresh-token",
		Expiry:       time.Now().Add(time.Minute),
	}
}
