package tokenmanager

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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

func TestKeyringRepository_Set(t *testing.T) {
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

	if err := repo.Set("user@domain.com", getDefaultToken()); err != nil {
		t.Errorf("error was not expected: err=%s", err)
	}
}

func TestKeyringRepository_Get(t *testing.T) {
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
	err = repo.Set("user@domain.com", want)
	if err != nil {
		t.Fatalf("error was not expected: err=%s", err)
	}

	t.Run("ShouldSuccess", func(t *testing.T) {
		got, err := repo.Get("user@domain.com")
		if err != nil {
			t.Errorf("error was not expected: err=%s", err)
		}

		if reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}
	})

	t.Run("ReturnErrNotFoundWhenTokenDoesNotExists", func(t *testing.T) {
		_, err := repo.Get("non-existent")
		if err != ErrTokenNotFound {
			t.Errorf("want: %s, got: %v", ErrTokenNotFound, err)
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

//type mockedPasswordReader struct {
//	value string
//}
//
//func (m *mockedPasswordReader) ReadPassword() (string, error) {
//	return m.value, nil
//}

func TestGetPassphraseFromEnvOrUserInputFn(t *testing.T) {
	//want := "foo"

	t.Run("Should return the passphrase from the user input", func(t *testing.T) {
		getPassphraseFromEnvOrUserInputFn := getPassphraseFromEnvOrUserInputFn()
		_, err := getPassphraseFromEnvOrUserInputFn("")

		// It should fail because this is not an interactive terminal.
		assert.Error(t, err)
	})

	t.Run("Should return the passphrase from the env var", func(t *testing.T) {
		if err := os.Setenv("GPHOTOS_CLI_TOKENSTORE_KEY", "This-key-comes-from-env-var"); err != nil {
			t.Fatalf("error was not expected at this stage: err=%s", err)
		}

		getPassphraseFromEnvOrUserInputFn := getPassphraseFromEnvOrUserInputFn()
		got, err := getPassphraseFromEnvOrUserInputFn("")

		assert.NoError(t, err)
		assert.Equal(t, "This-key-comes-from-env-var", got)

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
	return filepath.Join(os.TempDir(), fmt.Sprintf("gphotos-uploader-cli.%d", time.Now().UnixNano()))
}
