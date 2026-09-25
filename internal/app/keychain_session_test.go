package app

import (
	"errors"
	"strings"
	"testing"

	"github.com/zalando/go-keyring"
)

type fakeKeychain struct {
	value string
	err   error
}

func (f *fakeKeychain) Set(service, account, value string) error {
	if service != EvidenceKeychainService || account != EvidenceKeychainAccount {
		return errors.New("unexpected keychain location")
	}
	f.value = value
	return f.err
}

func (f fakeKeychain) Get(service, account string) (string, error) {
	if service != EvidenceKeychainService || account != EvidenceKeychainAccount {
		return "", errors.New("unexpected keychain location")
	}
	return f.value, f.err
}

func TestStoreEvidenceSession(t *testing.T) {
	keychain := &fakeKeychain{}
	if err := StoreEvidenceSession(keychain, " token "); err != nil {
		t.Fatalf("StoreEvidenceSession() error = %v", err)
	}
	if keychain.value != "token" {
		t.Fatalf("stored value = %q, want token", keychain.value)
	}
}

func TestReadEvidenceSession(t *testing.T) {
	got, err := ReadEvidenceSession(fakeKeychain{value: " token "})
	if err != nil {
		t.Fatalf("ReadEvidenceSession() error = %v", err)
	}
	if got != "token" {
		t.Fatalf("ReadEvidenceSession() = %q, want token", got)
	}
}

func TestReadEvidenceSessionFailsClosedWithoutLeakingValue(t *testing.T) {
	for _, test := range []struct {
		name string
		fake fakeKeychain
	}{
		{name: "missing", fake: fakeKeychain{err: keyring.ErrNotFound}},
		{name: "empty", fake: fakeKeychain{value: " "}},
		{name: "read failure", fake: fakeKeychain{err: errors.New("secret failure")}},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := ReadEvidenceSession(test.fake)
			if err == nil {
				t.Fatal("ReadEvidenceSession() error = nil")
			}
			if strings.Contains(err.Error(), "token") || strings.Contains(err.Error(), "secret") {
				t.Fatalf("error leaked credential detail: %v", err)
			}
		})
	}
}
