package tokenmanager_test

import (
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/tokenmanager"
)

const (
	ScenarioShouldSuccess = "should-success"
	ScenarioInvalidToken  = "invalid-token"
	ScenarioTokenNotFound = "token-not-found"
	ScenarioRefreshToken  = "refresh-token"
)

var valueInRepo oauth2.Token

func TestTokenManager_Get(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{"Should success", ScenarioShouldSuccess, nil},
		{"Should fail when token is not found", ScenarioTokenNotFound, tokenmanager.ErrTokenNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			want := getTestToken(tc.input)
			valueInRepo = want
			tm := tokenmanager.New(&mockedRepository{})

			got, err := tm.Get(tc.input)
			if tc.expectedErr != err {
				t.Errorf("err, want: %s, got: %s", tc.expectedErr, err)
			}
			if err == nil && !equalTokens(want, *got) {
				t.Errorf("token, want: %v, got: %v", want, got)
			}
		})
	}
}

func TestTokenManager_Put(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedErr error
	}{
		{"Should success", ScenarioShouldSuccess, nil},
		{"Should success besides refresh token is empty", ScenarioRefreshToken, nil},
		{"Should fail when token is invalid", ScenarioInvalidToken, tokenmanager.ErrInvalidToken},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			want := getTestToken(tc.input)
			mr := &mockedRepository{}
			tm := tokenmanager.New(mr)

			err := tm.Put(tc.input, &want)
			if tc.expectedErr != err {
				t.Errorf("err, want: %s, got: %s", tc.expectedErr, err)
			}

			if err == nil && !equalTokens(want, valueInRepo) {
				t.Errorf("token, want: %v, got: %v", want, valueInRepo)
			}
		})
	}
}

func TestTokenManager_Close(t *testing.T) {
	tm := tokenmanager.New(&mockedRepository{})

	t.Run("Should Success", func(t *testing.T) {
		if err := tm.Close(); err != nil {
			t.Errorf("error was not expected, err: %s", err)
		}
	})
}

type mockedRepository struct{}

func (mr mockedRepository) Get(key string) (*oauth2.Token, error) {
	if key == ScenarioShouldSuccess {
		token := valueInRepo
		return &token, nil
	}
	return nil, tokenmanager.ErrTokenNotFound
}

func (mr mockedRepository) Set(key string, token *oauth2.Token) error {
	valueInRepo = oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
	return nil
}

func (mr mockedRepository) Close() error {
	return nil
}

// getDefaultToken return a token to complete tests
func getTestToken(base string) oauth2.Token {
	token := oauth2.Token{
		AccessToken:  base + "-AccessToken",
		TokenType:    "Bearer",
		RefreshToken: base + "-RefreshToken",
		Expiry:       time.Now().Add(time.Hour),
	}

	switch base {
	case ScenarioInvalidToken:
		token.AccessToken = ""
	case ScenarioRefreshToken:
		token.RefreshToken = ""
	}
	return token
}

func equalTokens(want, in oauth2.Token) bool {
	if want.AccessToken == in.AccessToken &&
		want.RefreshToken == in.RefreshToken &&
		want.TokenType == in.TokenType &&
		want.Expiry == in.Expiry {
		return true
	}
	return false
}
