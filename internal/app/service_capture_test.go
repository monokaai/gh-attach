package app

import (
	"net/http"
	"testing"

	"github.com/sudosubin/gh-attach/internal/browserprovider"
	"github.com/sudosubin/gh-attach/internal/cookies"
)

func TestCaptureBrowserSessionTokenDoesNotRequireDotcomUser(t *testing.T) {
	t.Parallel()

	service := NewService(nil)
	service.providers = map[cookies.Browser]browserprovider.BrowserProvider{
		cookies.BrowserChrome: stubProvider{sessions: []browserprovider.BrowserSession{{
			Browser: string(cookies.BrowserChrome),
			Cookies: []*http.Cookie{{
				Name:   "user_session",
				Value:  "session-value",
				Domain: ".github.com",
				Path:   "/",
			}},
		}}},
	}

	got, err := service.CaptureBrowserSessionToken(t.Context(), "github.com", cookies.ResolveInput{
		Browser: cookies.BrowserChrome,
		Profile: "Default",
	}, false)
	if err != nil {
		t.Fatalf("CaptureBrowserSessionToken() error = %v", err)
	}
	if got != "session-value" {
		t.Fatalf("CaptureBrowserSessionToken() = %q", got)
	}
}

func TestCaptureBrowserSessionTokenRejectsSessionWithoutUserSession(t *testing.T) {
	t.Parallel()

	service := NewService(nil)
	service.providers = map[cookies.Browser]browserprovider.BrowserProvider{
		cookies.BrowserChrome: stubProvider{sessions: []browserprovider.BrowserSession{{
			Browser: string(cookies.BrowserChrome),
			Cookies: []*http.Cookie{{Name: "dotcom_user", Value: "monokaai", Domain: ".github.com", Path: "/"}},
		}}},
	}

	if _, err := service.CaptureBrowserSessionToken(t.Context(), "github.com", cookies.ResolveInput{Browser: cookies.BrowserChrome}, false); err == nil {
		t.Fatal("CaptureBrowserSessionToken() error = nil")
	}
}
