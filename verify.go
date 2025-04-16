package main

import (
	"fmt"
	"os"
	"time"

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
		fmt.Println("✅ Verified token claims:")

		for k, v := range claims {
			if k == "exp" {
				if expFloat, ok := v.(float64); ok {
					expTime := time.Unix(int64(expFloat), 0)
					diff := time.Until(expTime)
					status := ""
					if diff > 0 {
						status = fmt.Sprintf("in %s", diff.Round(time.Second))
					} else {
						status = fmt.Sprintf("%s ago", -diff.Round(time.Second))
					}
					fmt.Printf("  exp: %s (%s)\n", expTime.Format(time.RFC1123), status)
				} else {
					fmt.Printf("  exp: %v (invalid type)\n", v)
				}
			} else {
				fmt.Printf("  %s: %v\n", k, v)
			}
		}
	} else {
		return fmt.Errorf("invalid token")
	}

	return nil
}
