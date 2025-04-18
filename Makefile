# Makefile for jwt CLI tool

BINARY_NAME=jwt
DIST_DIR=bin
VERSION=$(shell git describe --tags --always)
BUILD_TIME=$(shell TZ=America/Chicago date)

LD_FLAGS=-ldflags "-X main.version=$(VERSION) -X 'main.buildTime=$(BUILD_TIME)'"

.PHONY: all local release clean

all: local

## 🔧 Build for local dev
local:
	@echo "🔨 Building for local GOARCH and GOOS..."
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 \
	go build $(LD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME) .

## ☁️ Build for AWS Linux 2 (static Linux binary)
release:
	@echo "📦 Last 3 release tags:"
	@git tag --sort=-creatordate | head -n 3 || echo "(no tags yet)"
	@echo ""

	@latest=$$(git tag --sort=-creatordate | head -n 1); \
	if [ -z "$$latest" ]; then \
		suggested="v0.1.0"; \
	else \
		ver=$$(echo $$latest | sed 's/^v//'); \
		major=$$(echo $$ver | cut -d. -f1); \
		minor=$$(echo $$ver | cut -d. -f2); \
		patch=$$(echo $$ver | cut -d. -f3); \
		patch=$$((patch + 1)); \
		suggested="v$${major}.$${minor}.$${patch}"; \
	fi; \
	echo "💡 Suggested next tag: $$suggested"; \
	read -p "Enter new release tag [default $$suggested]: " tag; \
	tag=$${tag:-$$suggested}; \
	if git tag | grep -q "^$$tag$$"; then \
		echo "❌ Tag '$$tag' already exists."; exit 1; \
	fi; \
	echo "🏷️  Creating and pushing tag '$$tag'..."; \
	git tag $$tag; \
	git push origin $$tag; \
	echo "🚀 Tag '$$tag' pushed. GitHub Actions will now build and release."

## 🧹 Clean build artifacts
clean:
	find $(DIST_DIR) -type f ! -name '.gitkeep' -delete
