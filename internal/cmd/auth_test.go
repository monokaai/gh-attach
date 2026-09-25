package cmd

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestAuthCaptureStoresWithoutWritingToken(t *testing.T) {
	previousCapture, previousStore := captureBrowserSession, storeEvidenceSession
	captureBrowserSession = func(browser, profile, path string) (string, error) { return "secret-token", nil }
	stored := ""
	storeEvidenceSession = func(value string) error { stored = value; return nil }
	t.Cleanup(func() { captureBrowserSession, storeEvidenceSession = previousCapture, previousStore })

	var output bytes.Buffer
	cmd := NewCmdAuthCapture()
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"--browser", "chrome"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if stored != "secret-token" {
		t.Fatalf("stored = %q", stored)
	}
	if strings.Contains(output.String(), "secret-token") {
		t.Fatalf("output leaked token: %q", output.String())
	}
}

func TestAuthCaptureDoesNotStoreWhenCaptureFails(t *testing.T) {
	previousCapture, previousStore := captureBrowserSession, storeEvidenceSession
	captureBrowserSession = func(string, string, string) (string, error) { return "", errors.New("no session") }
	storeEvidenceSession = func(string) error { t.Fatal("store was called"); return nil }
	t.Cleanup(func() { captureBrowserSession, storeEvidenceSession = previousCapture, previousStore })

	cmd := NewCmdAuthCapture()
	cmd.SetArgs([]string{"--browser", "chrome"})
	if err := cmd.Execute(); err == nil {
		t.Fatal("Execute() error = nil")
	}
}

func TestAuthImportStoresWithoutWritingToken(t *testing.T) {
	previousStore := storeEvidenceSession
	stored := ""
	storeEvidenceSession = func(value string) error { stored = value; return nil }
	t.Cleanup(func() { storeEvidenceSession = previousStore })

	path := t.TempDir() + "/cookies.json"
	if err := os.WriteFile(path, []byte(`[{"name":"user_session","value":"secret-token","domain":".github.com","path":"/"}]`), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	cmd := NewCmdAuthImport()
	cmd.SetOut(&output)
	cmd.SetArgs([]string{"--cookie-file", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if stored != "secret-token" {
		t.Fatalf("stored = %q", stored)
	}
	if strings.Contains(output.String(), "secret-token") {
		t.Fatalf("output leaked token: %q", output.String())
	}
}
