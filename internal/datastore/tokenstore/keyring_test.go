package tokenstore_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/99designs/keyring"
	"golang.org/x/oauth2"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/datastore/tokenstore"
)

var fixedStringPrompt keyring.PromptFunc = func(_ string) (string, error) {
	return "no more secrets", nil
}

func TestKeyringRepository_StoreToken(t *testing.T) {
	dir := tempDir()
	repo, err := tokenstore.NewKeyringRepository("file", &fixedStringPrompt, dir)
	if err != nil {
		t.Fatalf("error was not expected at this stage: err=%s", err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(err)
		}
	}()

	t.Run("ShouldReturnSuccess", func(t *testing.T) {
		if err := repo.StoreToken("user@domain.com", getDefaultToken()); err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}
	})

	t.Run("ShouldReturnErrInvalidTokenWhenTokenIsEmpty", func(t *testing.T) {
		token := &oauth2.Token{}
		if err := repo.StoreToken("user@domain.com", token); err != tokenstore.ErrInvalidToken {
			t.Errorf("want: %s, got: %v", tokenstore.ErrInvalidToken, err)
		}
	})
}

func TestKeyringRepository_RetrieveToken(t *testing.T) {
	dir := tempDir()
	repo, err := tokenstore.NewKeyringRepository("file", &fixedStringPrompt, dir)
	if err != nil {
		t.Fatalf("error was not expected at this stage: err=%s", err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(err)
		}
	}()

	want := getDefaultToken()
	err = repo.StoreToken("user@domain.com", want)
	if err != nil {
		t.Fatalf("error was not expected: err=%s", err)
	}

	t.Run("ShouldSuccess", func(t *testing.T) {
		got, err := repo.RetrieveToken("user@domain.com")
		if err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}

		if reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})

	t.Run("ReturnErrNotFoundWhenTokenDoesNotExists", func(t *testing.T) {
		_, err := repo.RetrieveToken("non-existent")
		if err != tokenstore.ErrNotFound {
			t.Errorf("want: %s, got: %v", tokenstore.ErrNotFound, err)
		}
	})
}

func TestKeyringRepository_CloseShouldSuccess(t *testing.T) {
	dir := tempDir()
	repo, err := tokenstore.NewKeyringRepository("file", &fixedStringPrompt, dir)
	if err != nil {
		t.Fatalf("error was not expected at this stage: err=%s", err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(err)
		}
	}()

	if err := repo.Close(); err != nil {
		t.Errorf("error was not expected: err=%s", err)
	}
}

// getDefaultToken return a token to complete tests
func getDefaultToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  "my-access-token",
		TokenType:    "my-token-type",
		RefreshToken: "my-refresh-token",
		Expiry:       time.Now().Add(time.Minute),
	}
}

func tempDir() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-cli.%d", time.Now().UnixNano()))
}
