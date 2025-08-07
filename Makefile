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

.PHONY: all build-api clean test format release

all: format test build-api

build-api:
	@echo "Building binaries for linux, windows, and macOS..."
	GOOS=linux   GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 \
		-ldflags "-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)'" $(CMD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe \
		-ldflags "-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)'" $(CMD_DIR)
	GOOS=darwin  GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-darwin-amd64 \
		-ldflags "-X 'main.Version=$(VERSION)' -X 'main.Commit=$(COMMIT)'" $(CMD_DIR)
	@echo "All builds complete. SHA256:"
	@echo "- linux:   $$(sha256sum $(BIN_DIR)/$(APP_NAME)-linux-amd64 | cut -d ' ' -f1)"
	@echo "- windows: $$(sha256sum $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe | cut -d ' ' -f1)"
	@echo "- darwin:  $$(sha256sum $(BIN_DIR)/$(APP_NAME)-darwin-amd64 | cut -d ' ' -f1)"

test:
	go test -v ./... \
		-ldflags "\
			-X 'main.Version=$(VERSION)' \
			-X 'main.Commit=$(COMMIT)'"

release: build-api
	@[ -n "$$$(SIGNING_KEY_ENV)" ] || (echo "FATAL: Environment variable $${SIGNING_KEY_ENV} is not set. Cannot sign release artifacts and manifest."; exit 1)

	@echo "Generating release manifest..."
	@SIGN_KEY_FILE=$$(mktemp); \
	echo "$$$(SIGNING_KEY_ENV)" > $$SIGN_KEY_FILE; \
	LIN_DIGEST=$$(sha256sum $(BIN_DIR)/$(APP_NAME)-linux-amd64 | cut -d ' ' -f1); \
	WIN_DIGEST=$$(sha256sum $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe | cut -d ' ' -f1); \
	MAC_DIGEST=$$(sha256sum $(BIN_DIR)/$(APP_NAME)-darwin-amd64 | cut -d ' ' -f1); \
	LIN_SIG=$$($(SIGN_CMD) $$SIGN_KEY_FILE $(BIN_DIR)/$(APP_NAME)-linux-amd64); \
	WIN_SIG=$$($(SIGN_CMD) $$SIGN_KEY_FILE $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe); \
	MAC_SIG=$$($(SIGN_CMD) $$SIGN_KEY_FILE $(BIN_DIR)/$(APP_NAME)-darwin-amd64); \
	UPDATED=$$(mktemp); \
	jq \
	  --arg version "$(VERSION)" \
	  --arg commit "$(COMMIT)" \
	  --arg pubkey "$$(printf %s "$(PUBLIC_KEY)")" \
	  --arg lin_digest "$$LIN_DIGEST" \
	  --arg win_digest "$$WIN_DIGEST" \
	  --arg mac_digest "$$MAC_DIGEST" \
	  --arg lin_sig "$$LIN_SIG" \
	  --arg win_sig "$$WIN_SIG" \
	  --arg mac_sig "$$MAC_SIG" \
	  -f scripts/merge_manifest.jq $(MANIFEST) > $$UPDATED; \
	mv $$UPDATED $(MANIFEST); \
	$(SIGN_CMD) $$SIGN_KEY_FILE $(MANIFEST) > $(MANIFEST).sig.base64; \
	rm -f $$SIGN_KEY_FILE; \
	echo "Manifest and signature updated: $(MANIFEST)"

clean:
	rm -rf \
		$(BIN_DIR)/$(APP_NAME)-linux-amd64 \
		$(BIN_DIR)/$(APP_NAME)-windows-amd64.exe \
		$(BIN_DIR)/$(APP_NAME)-darwin-amd64

format:
	go fmt ./...
