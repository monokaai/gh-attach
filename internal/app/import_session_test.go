package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadEvidenceSessionExportAcceptsArrayAndEnvelope(t *testing.T) {
	t.Parallel()

	for name, content := range map[string]string{
		"array":    `[{"name":"user_session","value":"session-value","domain":".github.com","path":"/"}]`,
		"envelope": `{"cookies":[{"name":"user_session","value":"session-value","domain":".github.com","path":"/"}]}`,
	} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "cookies.json")
			if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			got, err := ReadEvidenceSessionExport(path)
			if err != nil {
				t.Fatalf("ReadEvidenceSessionExport() error = %v", err)
			}
			if got != "session-value" {
				t.Fatalf("ReadEvidenceSessionExport() = %q", got)
			}
		})
	}
}

func TestReadEvidenceSessionExportRejectsNonGitHubSession(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "cookies.json")
	if err := os.WriteFile(path, []byte(`[{"name":"user_session","value":"session-value","domain":"example.com","path":"/"}]`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadEvidenceSessionExport(path); err == nil {
		t.Fatal("ReadEvidenceSessionExport() error = nil")
	}
}
