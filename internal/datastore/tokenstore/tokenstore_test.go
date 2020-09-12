package tokenstore

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/99designs/keyring"
	"golang.org/x/oauth2"
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

const secretsBackend string = "file"

var fixedStringPrompt keyring.PromptFunc = func(_ string) (string, error) {
	return "no more secrets", nil
}

func TestStoreToken(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))
	repo, err := NewKeyringRepository(secretsBackend, &fixedStringPrompt, dir)
	if err != nil {
		t.Errorf("error not expected at this stage: %v", err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(err)
		}
	}()
	s := NewService(repo)

	err = s.StoreToken(userEmail, getDefaultToken())
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}
}

// TestRetrieveToken tests getting a token from the keyring.
func TestRetrieveToken(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))
	repo, err := NewKeyringRepository(secretsBackend, &fixedStringPrompt, dir)
	if err != nil {
		t.Errorf("error not expected at this stage: %v", err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(err)
		}
	}()
	s := NewService(repo)

	expectedToken := getDefaultToken()
	err = s.StoreToken(userEmail, expectedToken)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	tk, err := s.RetrieveToken(userEmail)
	if err != nil {
		t.Errorf("Should not fail, got %s", err)
	}

	if tk.AccessToken != expectedToken.AccessToken {
		t.Errorf("Token doens't mismatch: expected %v, got %v", expectedToken, tk)
	}
}

// TestRetrieveExpiredToken tests getting an invalid (expired) token from the keyring.
func TestRetrieveExpiredToken(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))
	repo, err := NewKeyringRepository(secretsBackend, &fixedStringPrompt, dir)
	if err != nil {
		t.Errorf("error not expected at this stage: %v", err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(err)
		}
	}()
	s := NewService(repo)

	expectedToken := getDefaultToken()
	expectedToken.Expiry = time.Now().Add(-time.Minute)
	err = s.StoreToken(userEmail, expectedToken)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	tk, err := s.RetrieveToken(userEmail)
	if err != nil {
		t.Errorf("Should not fail, got %s", err)
	}

	if tk.AccessToken != expectedToken.AccessToken {
		t.Errorf("Token doens't mismatch: expected %v, got %v", expectedToken, tk)
	}
}

// TestRetrieveInvalidToken tests getting an invalid (empty AccessToken) token from the keyring.
func TestRetrieveInvalidToken(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))
	repo, err := NewKeyringRepository(secretsBackend, &fixedStringPrompt, dir)
	if err != nil {
		t.Errorf("error not expected at this stage: %v", err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(err)
		}
	}()
	s := NewService(repo)

	expectedToken := getDefaultToken()
	expectedToken.AccessToken = ""
	err = s.StoreToken(userEmail, expectedToken)
	if err != nil {
		t.Errorf("Should not fail, got: %s", err)
	}

	_, err = s.RetrieveToken(userEmail)
	if err != ErrInvalidToken {
		t.Errorf("Expected error ErrInvalidToken, got %s", err)
	}
}

// TestRetrieveNonExistingToken tests getting a token not in the keyring.
func TestRetrieveNonExistingToken(t *testing.T) {
	dir := filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-config.%d", time.Now().UnixNano()))
	repo, err := NewKeyringRepository(secretsBackend, &fixedStringPrompt, dir)
	if err != nil {
		t.Errorf("error not expected at this stage: %v", err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(err)
		}
	}()
	s := NewService(repo)

	_, err = s.RetrieveToken(userEmail + "fake")
	if err != ErrNotFound {
		t.Errorf("Expected error ErrNotFound, got %s", err)
	}
}
