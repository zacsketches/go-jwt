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

	if keyEnv != "" {
		envVal := os.Getenv(keyEnv)
		if envVal == "" {
			return fmt.Errorf("environment variable %s not set", keyEnv)
		}
		keyBytes, err := base64.StdEncoding.DecodeString(envVal)
		if err != nil {
			return fmt.Errorf("base64 decode error: %w", err)
		}
		privKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
		if err != nil {
			return fmt.Errorf("parsing private key from env: %w", err)
		}
	} else {
		keyData, err := os.ReadFile(keyPath)
		if err != nil {
			return fmt.Errorf("reading private key file: %w", err)
		}
		privKey, err = jwt.ParseRSAPrivateKeyFromPEM(keyData)
		if err != nil {
			return fmt.Errorf("parsing private key: %w", err)
		}
	}

	claims := jwt.MapClaims{
		"iss": issuer,
		"exp": time.Now().Add(time.Duration(expirySeconds) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(privKey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}

	fmt.Println(signedToken)
	return nil
}
