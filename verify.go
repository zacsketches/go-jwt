package main

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(publicKeyPath string, tokenString string) error {
	keyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("reading public key: %w", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return fmt.Errorf("parsing public key: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil {
		return fmt.Errorf("parsing token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Verified token claims:")
		for k, v := range claims {
			fmt.Printf("  %s: %v\n", k, v)
		}
	} else {
		return fmt.Errorf("invalid token")
	}

	return nil
}
