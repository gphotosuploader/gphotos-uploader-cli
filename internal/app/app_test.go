package app_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/mock"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/oauth"
)

const testAccount = "user@example.com"

// testConfig returns a valid configuration pointing to an existing source folder.
func testConfig(t *testing.T) config.Config {
	t.Helper()
	return config.Config{
		APIAppCredentials: config.APIAppCredentials{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
		},
		Account:            testAccount,
		SecretsBackendType: "file",
		Jobs: []config.FolderUploadJob{
			{SourceFolder: t.TempDir()},
		},
	}
}

// setupAppDir creates an application data dir containing the given configuration.
func setupAppDir(t *testing.T, cfg config.Config) string {
	t.Helper()
	t.Setenv("GPHOTOS_CLI_TOKENSTORE_KEY", "test-passphrase")

	dir := t.TempDir()
	b, err := json.Marshal(cfg)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, app.DefaultConfigFilename), b, 0600))
	return dir
}

// storeToken saves a token for the test account in the application token store.
func storeToken(t *testing.T, dir string, token *oauth2.Token) {
	t.Helper()
	a, err := app.StartServices(context.Background(), dir)
	require.NoError(t, err)
	require.NoError(t, a.TokenManager.Put(testAccount, token))
	require.NoError(t, a.Stop())
}

// fakeTokenEndpoint points the OAuth endpoint to a local test server.
func fakeTokenEndpoint(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	original := oauth.GoogleAuthEndpoint
	oauth.GoogleAuthEndpoint = oauth2.Endpoint{AuthURL: srv.URL + "/auth", TokenURL: srv.URL + "/token"}
	t.Cleanup(func() { oauth.GoogleAuthEndpoint = original })
}

func TestStartServices(t *testing.T) {
	t.Run("Should success", func(t *testing.T) {
		dir := setupAppDir(t, testConfig(t))

		a, err := app.StartServices(context.Background(), dir)
		require.NoError(t, err)

		assert.NotNil(t, a.FileTracker)
		assert.NotNil(t, a.TokenManager)
		assert.NotNil(t, a.UploadSessionTracker)
		assert.Equal(t, testAccount, a.Config.Account)
		assert.NoError(t, a.Stop())
	})

	t.Run("Should fail when configuration doesn't exist", func(t *testing.T) {
		_, err := app.StartServices(context.Background(), t.TempDir())
		assert.ErrorContains(t, err, "invalid configuration")
	})

	t.Run("Should fail when configuration is invalid", func(t *testing.T) {
		cfg := testConfig(t)
		cfg.Account = ""
		dir := setupAppDir(t, cfg)

		_, err := app.StartServices(context.Background(), dir)
		assert.ErrorContains(t, err, "invalid configuration")
	})

	t.Run("Should fail when file tracker fails to start", func(t *testing.T) {
		dir := setupAppDir(t, testConfig(t))
		// A regular file where the LevelDB folder is expected makes the repository fail.
		require.NoError(t, os.WriteFile(filepath.Join(dir, "uploaded_files"), nil, 0600))

		_, err := app.StartServices(context.Background(), dir)
		assert.ErrorContains(t, err, "file tracker could not be started")
	})

	t.Run("Should fail when token manager fails to start", func(t *testing.T) {
		if runtime.GOOS == "darwin" {
			t.Skip("keychain backend is available on macOS")
		}
		cfg := testConfig(t)
		cfg.SecretsBackendType = "keychain"
		dir := setupAppDir(t, cfg)

		_, err := app.StartServices(context.Background(), dir)
		assert.ErrorContains(t, err, "token manager could not be started")
	})

	t.Run("Should fail when uploads session tracker fails to start", func(t *testing.T) {
		dir := setupAppDir(t, testConfig(t))
		require.NoError(t, os.WriteFile(filepath.Join(dir, "ongoing_uploads"), nil, 0600))

		_, err := app.StartServices(context.Background(), dir)
		assert.ErrorContains(t, err, "uploads session tracker could not be started")
	})
}

func TestStart(t *testing.T) {
	t.Run("Should success with a valid token", func(t *testing.T) {
		dir := setupAppDir(t, testConfig(t))
		storeToken(t, dir, &oauth2.Token{
			AccessToken: "valid-access-token",
			Expiry:      time.Now().Add(time.Hour),
		})

		a, err := app.Start(context.Background(), dir)
		require.NoError(t, err)
		assert.NotNil(t, a.Client)
		assert.NoError(t, a.Stop())
	})

	t.Run("Should refresh and store an expired token", func(t *testing.T) {
		fakeTokenEndpoint(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"refreshed-access-token","token_type":"Bearer","expires_in":3600}`))
		})
		dir := setupAppDir(t, testConfig(t))
		storeToken(t, dir, &oauth2.Token{
			AccessToken:  "expired-access-token",
			RefreshToken: "refresh-token",
			Expiry:       time.Now().Add(-time.Hour),
		})

		a, err := app.Start(context.Background(), dir)
		require.NoError(t, err)
		t.Cleanup(func() { _ = a.Stop() })

		stored, err := a.TokenManager.Get(testAccount)
		require.NoError(t, err)
		assert.Equal(t, "refreshed-access-token", stored.AccessToken)
	})

	t.Run("Should fail when there is no token", func(t *testing.T) {
		dir := setupAppDir(t, testConfig(t))

		_, err := app.Start(context.Background(), dir)
		assert.ErrorContains(t, err, "unable to retrieve token")
	})

	t.Run("Should fail when token refresh fails", func(t *testing.T) {
		fakeTokenEndpoint(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
		})
		dir := setupAppDir(t, testConfig(t))
		storeToken(t, dir, &oauth2.Token{
			AccessToken:  "expired-access-token",
			RefreshToken: "revoked-refresh-token",
			Expiry:       time.Now().Add(-time.Hour),
		})

		_, err := app.Start(context.Background(), dir)
		assert.Error(t, err)
	})

	t.Run("Should fail when services fail to start", func(t *testing.T) {
		_, err := app.Start(context.Background(), t.TempDir())
		assert.ErrorContains(t, err, "invalid configuration")
	})
}

func TestApp_Stop(t *testing.T) {
	errClose := errors.New("close failed")

	testCases := []struct {
		name               string
		fileTrackerErr     error
		tokenManagerErr    error
		wantErr            error
		wantTokenMgrClosed bool
		wantSessionsClosed bool
	}{
		{"Should success", nil, nil, nil, true, true},
		{"Should fail when file tracker fails to close", errClose, nil, errClose, false, false},
		{"Should fail when token manager fails to close", nil, errClose, errClose, true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ft := &fakeFileTracker{closeErr: tc.fileTrackerErr}
			tm := &fakeTokenManager{closeErr: tc.tokenManagerErr}
			ust := &fakeUploadSessionTracker{}

			a := &app.App{
				FileTracker:          ft,
				TokenManager:         tm,
				UploadSessionTracker: ust,
				Logger:               &mock.Logger{},
			}

			err := a.Stop()
			assert.ErrorIs(t, err, tc.wantErr)
			assert.True(t, ft.closed)
			assert.Equal(t, tc.wantSessionsClosed, ust.closed)
			assert.Equal(t, tc.wantTokenMgrClosed, tm.closed)
		})
	}
}

func TestApp_CreateAppDataDir(t *testing.T) {
	t.Run("Should create the configuration file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		a, err := app.StartWithoutConfig(fs, "/app")
		require.NoError(t, err)
		assert.False(t, a.AppDataDirExists())

		filename, err := a.CreateAppDataDir()
		require.NoError(t, err)
		assert.Equal(t, filepath.Join("/app", app.DefaultConfigFilename), filename)
		assert.True(t, a.AppDataDirExists())
	})

	t.Run("Should remove previous application data", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		require.NoError(t, afero.WriteFile(fs, "/app/old-data", []byte("old"), 0600))

		a, err := app.StartWithoutConfig(fs, "/app")
		require.NoError(t, err)

		_, err = a.CreateAppDataDir()
		require.NoError(t, err)

		exists, err := afero.Exists(fs, "/app/old-data")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Should fail when the file system is read only", func(t *testing.T) {
		fs := afero.NewReadOnlyFs(afero.NewMemMapFs())
		a, err := app.StartWithoutConfig(fs, "/app")
		require.NoError(t, err)

		_, err = a.CreateAppDataDir()
		assert.Error(t, err)
		assert.False(t, a.AppDataDirExists())
	})
}

func TestApp_AuthenticateFromWeb_InvalidCredentials(t *testing.T) {
	a := &app.App{
		Config:       &config.Config{Account: testAccount},
		TokenManager: &fakeTokenManager{},
		Logger:       &mock.Logger{},
	}

	_, err := a.AuthenticateFromWeb(context.Background(), app.AuthenticationOptions{})
	assert.ErrorContains(t, err, "invalid config")
}

type fakeFileTracker struct {
	closeErr error
	closed   bool
}

func (f *fakeFileTracker) MarkAsUploaded(string) error   { return nil }
func (f *fakeFileTracker) IsUploaded(string) bool        { return false }
func (f *fakeFileTracker) UnmarkAsUploaded(string) error { return nil }
func (f *fakeFileTracker) Destroy() error                { return nil }
func (f *fakeFileTracker) Close() error {
	f.closed = true
	return f.closeErr
}

type fakeTokenManager struct {
	closeErr error
	closed   bool
}

func (f *fakeTokenManager) Put(string, *oauth2.Token) error   { return nil }
func (f *fakeTokenManager) Get(string) (*oauth2.Token, error) { return nil, errors.New("not found") }
func (f *fakeTokenManager) Close() error {
	f.closed = true
	return f.closeErr
}

type fakeUploadSessionTracker struct {
	closed bool
}

func (f *fakeUploadSessionTracker) Get(string) (string, bool) { return "", false }
func (f *fakeUploadSessionTracker) Set(string, string)        {}
func (f *fakeUploadSessionTracker) Delete(string)             {}
func (f *fakeUploadSessionTracker) Close()                    { f.closed = true }
