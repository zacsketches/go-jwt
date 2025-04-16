package main

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(keyPath string, keyEnv string, issuer string, expirySeconds int64) error {
	var privKey *rsa.PrivateKey
	var err error

	// Prefer explicitly passed --key-env
	if keyEnv != "" && keyEnv != "DEPLOY_SIGNING_KEY_B64" {
		envVal := os.Getenv(keyEnv)
		if envVal == "" {
			return fmt.Errorf("env var %s was specified but is not set", keyEnv)
		}
		privKey, err = parseKeyFromEnv(envVal)
		if err != nil {
			return err
		}
	} else {
		// Attempt file load first
		keyData, err := os.ReadFile(keyPath)
		if err == nil {
			privKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyData)
			if err != nil {
				return fmt.Errorf("failed to parse RSA private key from file: %w", err)
			}
		} else {
			// If file doesn't exist or fails to read, fallback to default env var
			envVal := os.Getenv("DEPLOY_SIGNING_KEY_B64")
			if envVal == "" {
				return fmt.Errorf("could not read key file (%s), and env var DEPLOY_SIGNING_KEY_B64 not set", keyPath)
			}
			privKey, err = parseKeyFromEnv(envVal)
			if err != nil {
				return err
			}
		}
	}

	claims := jwt.MapClaims{
		"iss": issuer,
		"exp": time.Now().Add(time.Duration(expirySeconds) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(privKey)
	if err != nil {
		return fmt.Errorf("failed to sign token: %w", err)
	}

	fmt.Println(signedToken)
	return nil
}

func parseKeyFromEnv(envVal string) (*rsa.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(envVal)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 from env: %w", err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key from env: %w", err)
	}
	return privKey, nil
}
