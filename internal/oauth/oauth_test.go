package oauth_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/oauth"
)

func validConfig() *oauth.Config {
	return &oauth.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
	}
}

func invalidConfigs() map[string]*oauth.Config {
	return map[string]*oauth.Config{
		"missing ClientID":     {ClientSecret: "client-secret"},
		"missing ClientSecret": {ClientID: "client-id"},
		"missing both":         {},
	}
}

// fakeTokenEndpoint starts a token endpoint and points oauth.GoogleAuthEndpoint to it.
func fakeTokenEndpoint(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	original := oauth.GoogleAuthEndpoint
	oauth.GoogleAuthEndpoint = oauth2.Endpoint{
		AuthURL:  srv.URL + "/auth",
		TokenURL: srv.URL + "/token",
	}
	t.Cleanup(func() { oauth.GoogleAuthEndpoint = original })
}

func TestGetToken_InvalidConfig(t *testing.T) {
	for name, cfg := range invalidConfigs() {
		t.Run(name, func(t *testing.T) {
			_, err := oauth.GetToken(context.Background(), cfg)
			assert.ErrorContains(t, err, "invalid config")
		})
	}
}

func TestGetToken_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cfg := validConfig()
	cfg.LocalServerBindAddress = []string{"127.0.0.1:0"}

	token, err := oauth.GetToken(ctx, cfg)
	assert.Error(t, err)
	assert.Nil(t, token)
}

func TestRefreshToken_InvalidConfig(t *testing.T) {
	for name, cfg := range invalidConfigs() {
		t.Run(name, func(t *testing.T) {
			_, err := oauth.RefreshToken(context.Background(), cfg, &oauth2.Token{})
			assert.ErrorContains(t, err, "invalid config")
		})
	}
}

func TestRefreshToken_NilToken(t *testing.T) {
	_, err := oauth.RefreshToken(context.Background(), validConfig(), nil)
	assert.ErrorIs(t, err, oauth.ErrTokenIsNil)
}

func TestRefreshToken_ValidTokenIsReturnedUnchanged(t *testing.T) {
	fakeTokenEndpoint(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("token endpoint should not be called for a valid token")
	})

	token := &oauth2.Token{
		AccessToken:  "valid-access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(time.Hour),
	}

	got, err := oauth.RefreshToken(context.Background(), validConfig(), token)
	require.NoError(t, err)
	assert.Same(t, token, got)
}

func TestRefreshToken_ExpiredTokenIsRefreshed(t *testing.T) {
	fakeTokenEndpoint(t, func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseForm())
		assert.Equal(t, "refresh_token", r.Form.Get("grant_type"))
		assert.Equal(t, "refresh-token", r.Form.Get("refresh_token"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"new-access-token","token_type":"Bearer","expires_in":3600}`))
	})

	token := &oauth2.Token{
		AccessToken:  "expired-access-token",
		RefreshToken: "refresh-token",
		Expiry:       time.Now().Add(-time.Hour),
	}

	got, err := oauth.RefreshToken(context.Background(), validConfig(), token)
	require.NoError(t, err)
	assert.Equal(t, "new-access-token", got.AccessToken)
	assert.True(t, got.Valid())
}

func TestRefreshToken_RefreshFails(t *testing.T) {
	fakeTokenEndpoint(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
	})

	token := &oauth2.Token{
		AccessToken:  "expired-access-token",
		RefreshToken: "revoked-refresh-token",
		Expiry:       time.Now().Add(-time.Hour),
	}

	_, err := oauth.RefreshToken(context.Background(), validConfig(), token)
	var retrieveErr *oauth2.RetrieveError
	assert.True(t, errors.As(err, &retrieveErr), "expected *oauth2.RetrieveError, got %v", err)
}

func TestClient_InvalidConfig(t *testing.T) {
	for name, cfg := range invalidConfigs() {
		t.Run(name, func(t *testing.T) {
			_, err := oauth.Client(context.Background(), cfg, &oauth2.Token{})
			assert.ErrorContains(t, err, "invalid config")
		})
	}
}

func TestClient_AuthorizesRequests(t *testing.T) {
	var gotAuthorization string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuthorization = r.Header.Get("Authorization")
	}))
	t.Cleanup(srv.Close)

	token := &oauth2.Token{
		AccessToken: "access-token",
		TokenType:   "Bearer",
		Expiry:      time.Now().Add(time.Hour),
	}

	client, err := oauth.Client(context.Background(), validConfig(), token)
	require.NoError(t, err)

	resp, err := client.Get(srv.URL)
	require.NoError(t, err)
	_ = resp.Body.Close()

	assert.Equal(t, "Bearer access-token", gotAuthorization)
}
