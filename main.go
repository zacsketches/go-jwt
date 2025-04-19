package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	version   = "dev"
	buildTime = "unknown"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "version":
		fmt.Printf("Version: %s\nBuilt:   %s\n", version, buildTime)

	case "sign":
		signCmd := flag.NewFlagSet("sign", flag.ExitOnError)
		keyPath := signCmd.String("key", "private.pem", "Path to RSA private key file")
		keyEnv := signCmd.String("key-env", "DEPLOY_SIGNING_KEY_B64", "Environment variable with base64-encoded RSA private key")
		issuer := signCmd.String("iss", "github-actions", "JWT issuer (default: github-actions)")
		expiry := signCmd.Int64("exp", 300, "Token expiry in seconds")

		signCmd.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage of jwt sign:\n")
			signCmd.PrintDefaults()
		}

		signCmd.Parse(os.Args[2:])
		err := SignToken(*keyPath, *keyEnv, *issuer, *expiry)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Sign error: %v\n", err)
			os.Exit(1)
		}

	case "verify":
		verifyCmd := flag.NewFlagSet("verify", flag.ExitOnError)
		publicKeyPath := verifyCmd.String("key", "public.pem", "Path to RSA public key file")
		token := verifyCmd.String("token", "", "JWT token to verify (required)")

		verifyCmd.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage of jwt verify:\n")
			verifyCmd.PrintDefaults()
		}

		verifyCmd.Parse(os.Args[2:])
		if *token == "" {
			fmt.Fprintln(os.Stderr, "Error: --token is required")
			verifyCmd.Usage()
			os.Exit(1)
		}
		err := VerifyToken(*publicKeyPath, *token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Verify error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`jwt - A minimal JWT sign/verify tool using RSA keys

Usage:
  jwt <command> [options]

Available commands:
  sign      Create a JWT using a private RSA key
  verify    Validate a JWT against a public RSA key
  version   Show the build version and timestamp

sign options:
  --key <file>          Path to RSA private key PEM file
  --env-key <var>       Name of env var holding base64-encoded private key
  --iss <string>        Issuer claim (default: "github-actions")
  --exp <int>           Expiration time in seconds (default: 300)
  
verify options:
  --key <file>          Path to RSA public key PEM file
  --env-key <var>       Name of env var holding base64-encoded public key
  --token <string>      JWT to verify

Examples:
  jwt sign --key private.pem --iss my-service --exp 10
  jwt verify --key public.pem --token <token>
  jwt version
`)
}
