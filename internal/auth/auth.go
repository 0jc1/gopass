package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
)

type Auth struct {
	pinHash     string
	currentPIN  string
}

func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) SetPIN(pin string) error {
	hash := sha256.Sum256([]byte(pin))
	a.pinHash = hex.EncodeToString(hash[:])
	a.currentPIN = pin
	
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	
	appDir := filepath.Join(configDir, "gopass")
	if err := os.MkdirAll(appDir, 0700); err != nil {
		return err
	}
	
	return os.WriteFile(filepath.Join(appDir, "pin.hash"), []byte(a.pinHash), 0600)
}

func (a *Auth) ValidatePIN(pin string) bool {
	hash := sha256.Sum256([]byte(pin))
	inputHash := hex.EncodeToString(hash[:])
	if a.pinHash == inputHash {
		a.currentPIN = pin
		return true
	}
	return false
}

func (a *Auth) LoadPINHash() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	
	pinHash, err := os.ReadFile(filepath.Join(configDir, "gopass", "pin.hash"))
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("PIN not set")
		}
		return err
	}
	
	a.pinHash = string(pinHash)
	return nil
}

func (a *Auth) IsPINSet() bool {
	return a.pinHash != ""
}

func (a *Auth) GetCurrentPIN() string {
	return a.currentPIN
}
