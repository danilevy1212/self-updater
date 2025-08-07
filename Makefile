APP_NAME := api
CMD_DIR := ./cmd/$(APP_NAME)
SIGN_CMD := go run ./cmd/sign
BIN_DIR := bin
ASSET_DIR := internal/assets
MANIFEST := $(ASSET_DIR)/release.json
PUBLIC_KEY := $(shell sed ':a;N;$$!ba;s/\n/\\n/g' internal/assets/public.pem)
SIGNING_KEY_ENV := SIGNING_KEY_PEM

VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)
BUILD_FLAGS=-ldflags "-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)'"

.PHONY: all build-api clean test format release

all: format test build-api

build-api:
	@echo "Building binaries for linux, windows, and macOS (amd64 + arm64)..."
	GOOS=linux   GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 $(BUILD_FLAGS) $(CMD_DIR)
	GOOS=linux   GOARCH=arm64 go build -o $(BIN_DIR)/$(APP_NAME)-linux-arm64 $(BUILD_FLAGS) $(CMD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe $(BUILD_FLAGS) $(CMD_DIR)
	GOOS=windows GOARCH=arm64 go build -o $(BIN_DIR)/$(APP_NAME)-windows-arm64.exe $(BUILD_FLAGS) $(CMD_DIR)
	GOOS=darwin  GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-darwin-amd64 $(BUILD_FLAGS) $(CMD_DIR)
	GOOS=darwin  GOARCH=arm64 go build -o $(BIN_DIR)/$(APP_NAME)-darwin-arm64 $(BUILD_FLAGS) $(CMD_DIR)
	@echo "All builds complete. SHA256:"
	@echo "- linux-amd64:   $$(sha256sum $(BIN_DIR)/$(APP_NAME)-linux-amd64 | cut -d ' ' -f1)"
	@echo "- linux-arm64:   $$(sha256sum $(BIN_DIR)/$(APP_NAME)-linux-arm64 | cut -d ' ' -f1)"
	@echo "- windows-amd64: $$(sha256sum $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe | cut -d ' ' -f1)"
	@echo "- windows-arm64: $$(sha256sum $(BIN_DIR)/$(APP_NAME)-windows-arm64.exe | cut -d ' ' -f1)"
	@echo "- darwin-amd64:  $$(sha256sum $(BIN_DIR)/$(APP_NAME)-darwin-amd64 | cut -d ' ' -f1)"
	@echo "- darwin-arm64:  $$(sha256sum $(BIN_DIR)/$(APP_NAME)-darwin-arm64 | cut -d ' ' -f1)"

test:
	go test -v ./... \
		-ldflags "\
			-X 'main.Version=$(VERSION)' \
			-X 'main.Commit=$(COMMIT)'"

release: build-api
	@sh -c 'VERSION="$(VERSION)" COMMIT="$(COMMIT)" PUBLIC_KEY="$(PUBLIC_KEY)" SIGNING_KEY_PEM="$${SIGNING_KEY_PEM}" ./scripts/release.sh'

clean:
	rm -rf \
		$(BIN_DIR)/$(APP_NAME)-linux-amd64 \
		$(BIN_DIR)/$(APP_NAME)-linux-arm64 \
		$(BIN_DIR)/$(APP_NAME)-windows-amd64.exe \
		$(BIN_DIR)/$(APP_NAME)-windows-arm64.exe \
		$(BIN_DIR)/$(APP_NAME)-darwin-amd64 \
		$(BIN_DIR)/$(APP_NAME)-darwin-arm64

format:
	go fmt ./...
