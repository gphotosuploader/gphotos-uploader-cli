package oauth

import (
	"net/url"
	"strings"
	"testing"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func TestLeastPrivilegeScopes(t *testing.T) {
	c := &Config{ClientID: "id", ClientSecret: "secret"}
	if err := c.validateAndSetDefaults(); err != nil {
		t.Fatal(err)
	}
	want := []string{PhotosLibraryAppendOnlyScope}
	if len(c.oAuth2Config.Scopes) != 1 || c.oAuth2Config.Scopes[0] != want[0] {
		t.Fatalf("OAuth scopes = %v, want append-only", c.oAuth2Config.Scopes)
	}
}

func TestOfflinePKCEAuthorizationParameters(t *testing.T) {
	c := &Config{ClientID: "id", ClientSecret: "secret"}
	GoogleAuthEndpoint = oauth2.Endpoint{AuthURL: "https://example.test/auth", TokenURL: "https://example.test/token"}
	t.Cleanup(func() { GoogleAuthEndpoint = google.Endpoint })
	if err := c.validateAndSetDefaults(); err != nil {
		t.Fatal(err)
	}
	verifier := oauth2.GenerateVerifier()
	web := c.webConfig(make(chan string, 1), verifier)
	raw := web.OAuth2Config.AuthCodeURL("state", web.AuthCodeOptions...)
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	if q.Get("access_type") != "offline" {
		t.Fatalf("access_type = %q", q.Get("access_type"))
	}
	if q.Get("prompt") != "consent" {
		t.Fatalf("prompt = %q", q.Get("prompt"))
	}
	if q.Get("code_challenge_method") != "S256" || q.Get("code_challenge") == "" {
		t.Fatalf("PKCE parameters missing: %s", raw)
	}
	if got := strings.Fields(q.Get("scope")); len(got) != 1 || got[0] != PhotosLibraryAppendOnlyScope {
		t.Fatalf("scope = %v", got)
	}
	if len(web.TokenRequestOptions) != 1 {
		t.Fatalf("token request options = %d, want PKCE verifier", len(web.TokenRequestOptions))
	}
}
