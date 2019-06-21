package tokenstore

import (
	"golang.org/x/oauth2"
	"testing"
	"time"
)

// getDefaultToken return a token to complete tests
func getDefaultToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  "my-access-token",
		TokenType:    "my-token-type",
		RefreshToken: "my-refresh-token",
		Expiry:       time.Now().Add(time.Minute),
	}
}

const userEmail string = "user@domain.com"

// TestStoreToken tests setting a token in the keyring.
func TestStoreToken(t *testing.T) {
	err := StoreToken(userEmail, getDefaultToken())
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}
}

// TestRetrieveToken tests getting a token from the keyring.
func TestRetrieveToken(t *testing.T) {
	expectedToken := getDefaultToken()
	err := StoreToken(userEmail, expectedToken)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	tk, err := RetrieveToken(userEmail)
	if err != nil {
		t.Errorf("Should not fail, got %s", err)
	}

	if tk.AccessToken != expectedToken.AccessToken {
		t.Errorf("Token doens't mismatch: expected %v, got %v", expectedToken, tk)
	}
}

// TestRetrieveExpiredToken tests getting an invalid (expired) token from the keyring.
func TestRetrieveExpiredToken(t *testing.T) {
	expectedToken := getDefaultToken()
	expectedToken.Expiry = time.Now().Add(-time.Minute)
	err := StoreToken(userEmail, expectedToken)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	_, err = RetrieveToken(userEmail)
	if err != ErrInvalidToken {
		t.Errorf("Expected error ErrInvalidToken, got %s", err)
	}
}

// TestRetrieveInvalidToken tests getting an invalid (empty AccessToken) token from the keyring.
func TestRetrieveInvalidToken(t *testing.T) {
	expectedToken := getDefaultToken()
	expectedToken.AccessToken = ""
	err := StoreToken(userEmail, expectedToken)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	_, err = RetrieveToken(userEmail)
	if err != ErrInvalidToken {
		t.Errorf("Expected error ErrInvalidToken, got %s", err)
	}
}

// TestRetrieveNonExistingToken tests getting a token not in the keyring.
func TestRetrieveNonExistingToken(t *testing.T) {
	_, err := RetrieveToken(userEmail + "fake")
	if err != ErrNotFound {
		t.Errorf("Expected error ErrNotFound, got %s", err)
	}
}
