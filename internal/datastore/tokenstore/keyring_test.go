package tokenstore

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/99designs/keyring"
	"golang.org/x/oauth2"
)

var fixedStringPrompt keyring.PromptFunc = func(_ string) (string, error) {
	return "no more secrets", nil
}

func TestKeyringRepository_StoreToken(t *testing.T) {
	dir := tempDir()
	repo, err := NewKeyringRepository("file", &fixedStringPrompt, dir)
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
		if err := repo.StoreToken("user@domain.com", token); err != ErrInvalidToken {
			t.Errorf("want: %s, got: %v", ErrInvalidToken, err)
		}
	})

	t.Run("ShouldStoreOldRefreshTokenIfItIsNotProvided", func(t *testing.T) {
		oldToken := getDefaultToken()
		if err := repo.StoreToken("user@domain.com", oldToken); err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}

		// newToken is not defining the RefreshToken.
		newToken := &oauth2.Token{
			AccessToken: "my-new-access-token",
		}
		if err := repo.StoreToken("user@domain.com", newToken); err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}

		got, err := repo.RetrieveToken("user@domain.com")
		if err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}
		if got.RefreshToken != oldToken.RefreshToken {
			t.Errorf("want: %s, got: %s", oldToken.RefreshToken, got.RefreshToken)
		}
	})
}

func TestKeyringRepository_RetrieveToken(t *testing.T) {
	dir := tempDir()
	repo, err := NewKeyringRepository("file", &fixedStringPrompt, dir)
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
		if err != ErrNotFound {
			t.Errorf("want: %s, got: %v", ErrNotFound, err)
		}
	})
}

func TestKeyringRepository_Close(t *testing.T) {
	dir := tempDir()
	repo, err := NewKeyringRepository("file", &fixedStringPrompt, dir)
	if err != nil {
		t.Fatalf("error was not expected at this stage: err=%s", err)
	}
	defer func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(err)
		}
	}()

	t.Run("ShouldSuccess", func(t *testing.T) {
		if err := repo.Close(); err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}
	})
}

type mockedPasswordReader struct {
	value string
}

func (m *mockedPasswordReader) ReadPassword() (string, error) {
	return m.value, nil
}

func TestPromptFn(t *testing.T) {
	want := "foo"

	t.Run("ReturnKeyFromTerminal", func(t *testing.T) {
		promptFn := promptFn(&mockedPasswordReader{value: want})
		got, err := promptFn("")
		if err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}
		if got != want {
			t.Errorf("want: %s, got: %s", want, got)
		}
	})

	t.Run("ReturnKeyFromEnv", func(t *testing.T) {
		if err := os.Setenv("GPHOTOS_CLI_TOKENSTORE_KEY", want); err != nil {
			t.Fatalf("error was not expected at this stage: err=%s", err)
		}

		promptFn := promptFn(&mockedPasswordReader{value: "dummy"})
		got, err := promptFn("")
		if err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}
		if got != want {
			t.Errorf("want: %s, got: %s", want, got)
		}
	})
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
