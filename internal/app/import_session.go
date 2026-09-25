package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/sudosubin/gh-attach/internal/cookies"
)

type exportedCookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type exportedCookieEnvelope struct {
	Cookies []exportedCookie `json:"cookies"`
}

// ReadEvidenceSessionExport reads a local browser cookie export and extracts
// only the GitHub user_session value. The caller must store the value in the
// Keychain and must never print or persist the source file's contents.
func ReadEvidenceSessionExport(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read cookie export: %w", err)
	}

	var list []exportedCookie
	if err := json.Unmarshal(data, &list); err != nil {
		var envelope exportedCookieEnvelope
		if envelopeErr := json.Unmarshal(data, &envelope); envelopeErr != nil {
			return "", errors.New("cookie export must be a JSON cookie array or an object with a cookies array")
		}
		list = envelope.Cookies
	}

	httpCookies := make([]*http.Cookie, 0, len(list))
	for _, cookie := range list {
		httpCookies = append(httpCookies, &http.Cookie{
			Name:   cookie.Name,
			Value:  cookie.Value,
			Domain: cookie.Domain,
			Path:   cookie.Path,
		})
	}
	values := cookies.ValuesForHost(httpCookies, "user_session", "github.com")
	if len(values) == 0 {
		return "", errors.New("cookie export does not contain a GitHub user_session")
	}
	return values[0], nil
}
