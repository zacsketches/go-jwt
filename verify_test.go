package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKeyPath       = "test-keys/private.pem"
	correctPublicKeyPath = "test-keys/public.pem"
	wrongPublicKeyPath   = "test-keys/other-pub.pem"
	missingKeyPath       = "test-keys/does-not-exist.pem"
)

// Generate a valid signed token for use in all tests
func generateValidToken(t *testing.T) string {
	t.Helper()

	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		t.Fatalf("Failed to read private key: %v", err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		t.Fatalf("Failed to parse private key: %v", err)
	}

	claims := jwt.MapClaims{
		"iss": "test-suite",
		"exp": time.Now().Add(5 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	return tokenString
}

func TestVerifyToken_CorrectPublicKey(t *testing.T) {
	token := generateValidToken(t)

	err := VerifyToken(correctPublicKeyPath, token)
	if err != nil {
		t.Errorf("Expected token to verify with correct public key, but got error: %v", err)
	}
}

func TestVerifyToken_MissingKeyFile(t *testing.T) {
	token := generateValidToken(t)

	err := VerifyToken(missingKeyPath, token)
	if err == nil {
		t.Error("Expected error due to missing public key file, but got none")
	} else if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("Expected missing file error, got: %v", err)
	}
}

func TestVerifyToken_WrongPublicKey(t *testing.T) {
	token := generateValidToken(t)

	// Make sure the wrong key file exists
	if _, err := os.Stat(wrongPublicKeyPath); os.IsNotExist(err) {
		t.Fatalf("Expected key file %s to exist, but got: %v", wrongPublicKeyPath, err)
	}

	err := VerifyToken(wrongPublicKeyPath, token)
	if err == nil {
		t.Error("Expected verification to fail with wrong public key, but it succeeded")
	} else if !strings.Contains(err.Error(), "crypto/rsa: verification error") {
		t.Errorf("Expected signature verification error, got: %v", err)
	}
}
