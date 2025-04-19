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
	@echo "📦 Building for AWS Linux (static)..."
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go build $(LD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux .

## 🧹 Clean build artifacts
clean:
	find $(DIST_DIR) -type f ! -name '.gitkeep' -delete
