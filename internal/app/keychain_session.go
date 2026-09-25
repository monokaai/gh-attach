package app

import (
	"errors"
	"strings"

	"github.com/zalando/go-keyring"
)

const (
	EvidenceKeychainService = "com.monokaai.gh-attach.evidence"
	EvidenceKeychainAccount = "github.com"
)

type KeychainReader interface {
	Get(service, account string) (string, error)
}

type KeychainWriter interface {
	Set(service, account, value string) error
}

type systemKeychain struct{}

func (systemKeychain) Get(service, account string) (string, error) {
	return keyring.Get(service, account)
}

func (systemKeychain) Set(service, account, value string) error {
	return keyring.Set(service, account, value)
}

func ReadEvidenceSession(reader KeychainReader) (string, error) {
	value, err := reader.Get(EvidenceKeychainService, EvidenceKeychainAccount)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", errors.New("headless evidence session is not registered; re-authenticate while macOS is unlocked")
		}
		return "", errors.New("headless evidence session is unavailable; re-authenticate while macOS is unlocked")
	}
	if strings.TrimSpace(value) == "" {
		return "", errors.New("headless evidence session is empty; re-authenticate while macOS is unlocked")
	}
	return strings.TrimSpace(value), nil
}

func ReadSystemEvidenceSession() (string, error) {
	return ReadEvidenceSession(systemKeychain{})
}

func StoreEvidenceSession(writer KeychainWriter, value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return errors.New("headless evidence session is empty")
	}
	if err := writer.Set(EvidenceKeychainService, EvidenceKeychainAccount, value); err != nil {
		return errors.New("store headless evidence session failed")
	}
	return nil
}

func StoreSystemEvidenceSession(value string) error {
	return StoreEvidenceSession(systemKeychain{}, value)
}
