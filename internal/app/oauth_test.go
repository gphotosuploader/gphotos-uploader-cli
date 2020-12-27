package app_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/oauth2"

	"github.com/gphotosuploader/gphotos-uploader-cli/internal/app"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/config"
	"github.com/gphotosuploader/gphotos-uploader-cli/internal/log"
)

const (
	ScenarioShouldSuccess          = "should-success"
	ScenarioFailedOAuth2           = "should-fail-oauth2"
	ScenarioSuccessfulExpiredToken = "successful-expired-token"
	ScenarioFailedExpiredToken     = "failed-expired-token"
)

func TestApp_NewOAuth2Client(t *testing.T) {
	testCases := []struct {
		name          string
		scenario      string
		isErrExpected bool
	}{
		{"Should success", ScenarioShouldSuccess, false},
		{"Should success if OAuth2 token is refreshed", ScenarioSuccessfulExpiredToken, false},
		{"Should fail if OAuth2 token is not refreshed", ScenarioFailedExpiredToken, true},
		{"Should fail if OAuth2 token fails", ScenarioFailedOAuth2, true},
	}

	srv := NewMockedGoogleAuthServer()
	defer srv.Close()

	// Set Google endpoints pointing to our test service.
	app.GoogleAuthEndpoint = oauth2.Endpoint{
		AuthURL:   srv.baseURL + "/auth",
		TokenURL:  srv.baseURL + "/token",
		AuthStyle: oauth2.AuthStyleInParams,
	}

	// Set non-interactive function to get authorization code.
	app.AskForAuthCodeFn = func(r io.Reader, url string) (string, error) {
		return url, nil
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testApp := newTestApp(tc.scenario)

			_, err := testApp.NewOAuth2Client(context.Background())

			assertExpectedError(t, tc.isErrExpected, err)
		})
	}
}

func newTestApp(scenario string) *app.App {
	var tokenManagerValue *oauth2.Token
	var tokenManagerErr error

	switch scenario {
	case ScenarioSuccessfulExpiredToken:
		tokenManagerValue = &oauth2.Token{
			RefreshToken: ScenarioSuccessfulExpiredToken,
		}
		tokenManagerErr = nil
	case ScenarioFailedExpiredToken:
		tokenManagerValue = &oauth2.Token{
			RefreshToken: ScenarioFailedExpiredToken,
		}
		tokenManagerErr = nil
	default:
		tokenManagerErr = errors.New("error-in-token-manager")
	}
	testApp := &app.App{
		TokenManager: &mockedTokenManager{
			RetrieveTokenFn: func(email string) (*oauth2.Token, error) {
				return tokenManagerValue, tokenManagerErr
			},
			StoreTokenFn: func(email string, token *oauth2.Token) error {
				return nil
			},
		},
		Logger: log.Discard,
		Config: &config.Config{
			APIAppCredentials: config.APIAppCredentials{
				ClientID: scenario,
			},
		},
	}
	return testApp
}

// mockedGoogleAuthServer mocks the Google authorization service.
type mockedGoogleAuthServer struct {
	server  *httptest.Server
	baseURL string
}

func NewMockedGoogleAuthServer() *mockedGoogleAuthServer {
	ms := &mockedGoogleAuthServer{}
	mux := http.NewServeMux()
	ms.server = httptest.NewServer(mux)
	ms.baseURL = ms.server.URL
	mux.HandleFunc("/token", ms.handleToken)
	return ms
}

func (ms mockedGoogleAuthServer) Close() {
	ms.server.Close()
}

func (ms mockedGoogleAuthServer) handleToken(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	clientId := getValueFromBody(body, "client_id")

	switch {
	case clientId == "" || clientId == ScenarioFailedOAuth2 || clientId == ScenarioFailedExpiredToken:
		w.WriteHeader(http.StatusInternalServerError)
	case clientId == ScenarioSuccessfulExpiredToken:
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token": "my-access-token", "scope": "user", "token_type": "bearer", "expires_in": 86400}`))
	default:
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token": "my-access-token", "scope": "user", "token_type": "bearer", "expires_in": 86400}`))
	}
}

// getValueFromBody returns the value for the specified key.
func getValueFromBody(body []byte, key string) string {
	ss := strings.Split(string(body), "&")
	for _, s := range ss {
		ss := strings.SplitN(s, "=", 2)
		if ss[0] == key {
			return ss[1]
		}
	}
	return ""
}

type mockedTokenManager struct {
	StoreTokenFn    func(email string, token *oauth2.Token) error
	RetrieveTokenFn func(email string) (*oauth2.Token, error)
	CloseFn         func() error
}

func (m mockedTokenManager) StoreToken(email string, token *oauth2.Token) error {
	return m.StoreTokenFn(email, token)
}

func (m mockedTokenManager) RetrieveToken(email string) (*oauth2.Token, error) {
	return m.RetrieveTokenFn(email)
}

func (m mockedTokenManager) Close() error {
	return m.CloseFn()
}

func assertExpectedError(t *testing.T, errExpected bool, err error, ) {
	if errExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !errExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}
