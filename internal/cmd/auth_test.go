package cmd

import (
	"bytes"
	"errors"
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

	if err := NewCmdAuthCapture().Execute(); err == nil {
		t.Fatal("Execute() error = nil")
	}
}
