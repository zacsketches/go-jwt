# go-jwt
Simple Go cli tool for signing and decoding jwt tokens.

In production this tool is called in github actions and creates a signed jwt to send with a webhook. The tool is opinionated about its defaults to make usage in Github actions very simple.

This tool difers from other tools on github because it statically compiles the required C library for `github.com/golang-jwt/jwt/v5` into the release, and avoids an error that often arises when a wrapper to a C library, like `golang-jwt/jwt/v4`, has as specific library dependency that is not installed on a cloud VM.  For example, the error that inspired the creation of this tool was this error when running a different jwt cli on AWS.
```
 /lib64/libc.so.6: version `GLIBC_2.34' not found (required by github.com/golang-jwt/jwt/v5)
```
This error presented on AWS Linux 2 in 2025 when the default EC2 is pre-installed with  `GLIBC_2.26`, which misses the dependency requirement of the CGO embedded library and crashed the other jwt cli tool.

Because it is designed to run in Github Actions and in the cloud `go-jwt` allows for testing on a cloud instance with high confidence that the webhook produced by Github Actions in production will look exactly like a local jwt produced on the instance to test the webhook handler.

#### Default usage is to set up for an environment variable called DEPLOY_SIGNING_KEY_B64 and relies on the default `iss` and `exp` times.
`jwt sign`

#### If you explicitly provide a file, it won’t look for the env var
`jwt sign --key test-keys/private.pem`

#### If you don’t pass --key and don't have an --env-key set, it will warn clearly
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
Defaults set `--key-env` to `DEPLOY_SIGNING_KEY_B64` so for simplest use copy the output of Step 1 into a Github Secret named `DEPLOY_SIGNING_KEY_B64`.

The `-iss` claim defaults to `"github-actions"` if not overridden.
The `-exp` claim defaults to `300` seconds if not overriden.

#### 3. In Github action:

# ⚠️ Update this section to include pulling the released version
```
- name: Sign JWT
  run: |
    export DEPLOY_SIGNING_KEY_B64="${{ secrets.DEPLOY_SIGNING_KEY_B64 }}"
    go run main.go sign 
```

## Local testing use
Primary access is provided through releases.  If you are going to build from source then you will need to include the build flags that support the `jwt version` command and statically compile dependencies as discussed above. Releases are built with the following flags:
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o bin/jwt . \
    -ldflags "-X main.version=$(git describe --tags --always) -X 'main.buildTime=$(TZ=America/Chicago date)'"
```

Outpus above reflect my preferred  convenience of `direnv` to add the `/go-jwt/bin` folder to my path dynamically whenever I enter the `/go-jwt` root of cloned Github project. Then default usage becomes
```
export DEPLOY_SIGNING_KEY_B64=$(base64 < <project root>/test-keys/private.pem)
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

