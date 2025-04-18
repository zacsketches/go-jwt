# go-jwt
Simple Go cli tool for signing and decoding jwt tokens.

In production this tool is called in github actions and creates a signed jwt to send with a webhook. The tool is opinionated about it's defaults to make usage in Github actions very simple.

This tool difers from other tools on github because it statically compiles the required C library for `github.com/golang-jwt/jwt/v5` into the release, and avoids this error often occurred when a wrapper to a C library like `golang-jwt` has as specific library dependency that is not installed on a cloud VM.  The error that inspired the creation of the tool was
```
 /lib64/libc.so.6: version `GLIBC_2.34' not found (required by github.com/golang-jwt/jwt/v5)
```
This error presented on AWS Linux 2 in 2025 when the default instance is at `GLIBC_2.26`.

In testing this tool is built to load onto a cloud instance for a webhook triggered deployment manager service because it generates tokens for local testing and verfication of the handler.

#### Normal usage is to set an environment variable called DEPLOY_SIGNING_KEY_B64 and rely on the default `iss` and `exp` times.
`jwt sign`

#### If you explicitly provide a file, it won’t look for the env var
`jwt sign --key test-keys/private.pem`

#### If you don’t pass --key and rely on the env var, it will warn clearly
`jwt sign`
`Sign error: could not read key file (private.pem), and env var DEPLOY_SIGNING_KEY_B64 not set`

## CI Usage
Primary usage in Github actions as follows
#### 1. Convert a local private key into a base 64 string.
This encoding is important to take out the newline characters in a typical multiline `*.pem` file that is generated when you create the key. The long single line of base 64 encoded data is much better suited to store and retrieve as an environment variable.
```
base64 /path/to/private.pem
```
#### 2. Create a Github Secret
Defaults set `-key-env` to `DEPLOY_SIGNING_KEY_B64` so for simplest use copy the output of Step 1 into a Github Secret named `DEPLOY_SIGNING_KEY_B64`.

The `-iss` claim defaults to `"github-actions"` if not overridden.
The `-exp` claim defaults to `300` seconds if not overriden.

#### 3. In Github action:
```
- name: Sign JWT
  run: |
    export DEPLOY_SIGNING_KEY_B64="${{ secrets.DEPLOY_SIGNING_KEY_B64 }}"
    go run main.go sign 
```

## Local testing use
Usage assumes something along the lines of `go build -o /bin/jwt`, and for convenience I use `direnv` to add the `.../go-jwt/bin` folder to my path dynamically whenever I enter the `/go-jwt` project foler.
```
export DEPLOY_SIGNING_KEY_B64=$(base64 < private.pem)
jwt sign
```

## Useful command line foo

```
openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in private.pem -out public.pem
```

Also genearate `other-priv.pem` and `other-pub.pem` so you can verify what it looks like when a token fails verification.
```
openssl genpkey -algorithm RSA -out other-priv.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in private.pem -out other-pub.pem
```

