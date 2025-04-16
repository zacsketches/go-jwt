# go-jwt
Simple Go cli tool for signing **and decoding** jwt tokens.

This project was built to create a tool I can pull into github actions and create a signed jwt to send with a webhook. The tool is opinionated about it's defaults to make usage in Github actions very simple.

Most importantly, I can also load this tool onto the cloud instance for a webhook triggered deployment manager service and generate tokens for local testing and verfication of the handler.

## CI Usage
Primary usage in Github actions as follows
#### 1. Store private.pem as a secret
```
base64 private.pem
```
#### 2. Create a Github Secret
Defaults set `--key-env` to `DEPLOY_SIGNING_KEY_B64` so for simplest use copy the output of Step 1 into a Github Secret named `DEPLOY_SIGNING_KEY_B64`.

The `--iss` claim defaults to `"github-actions"` if not overridden.

#### 3. In Github action:
```
- name: Sign JWT
  run: |
    export DEPLOY_SIGNING_KEY_B64="${{ secrets.DEPLOY_SIGNING_KEY_B64 }}"
    go run main.go sign 
```

## Local testing use
Usage assumes something along the lines of `go build -o jwt`.
```
export DEPLOY_SIGNING_KEY_B64=$(base64 < private.pem)
jwt sign
```

## Useful command line foo

```
openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in private.pem -out public.pem
```

