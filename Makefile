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
	@[ -n "$$$(SIGNING_KEY_ENV)" ] || (echo "FATAL: Environment variable $${SIGNING_KEY_ENV} is not set. Cannot sign release artifacts and manifest."; exit 1)

	@echo "Generating release manifest..."
	@SIGN_KEY_FILE=$$(mktemp); \
	echo "$$$(SIGNING_KEY_ENV)" > $$SIGN_KEY_FILE; \
	LIN_AMD_DIGEST=$$(sha256sum $(BIN_DIR)/$(APP_NAME)-linux-amd64 | cut -d ' ' -f1); \
	LIN_ARM_DIGEST=$$(sha256sum $(BIN_DIR)/$(APP_NAME)-linux-arm64 | cut -d ' ' -f1); \
	WIN_AMD_DIGEST=$$(sha256sum $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe | cut -d ' ' -f1); \
	WIN_ARM_DIGEST=$$(sha256sum $(BIN_DIR)/$(APP_NAME)-windows-arm64.exe | cut -d ' ' -f1); \
	MAC_AMD_DIGEST=$$(sha256sum $(BIN_DIR)/$(APP_NAME)-darwin-amd64 | cut -d ' ' -f1); \
	MAC_ARM_DIGEST=$$(sha256sum $(BIN_DIR)/$(APP_NAME)-darwin-arm64 | cut -d ' ' -f1); \
	LIN_AMD_SIG=$$($(SIGN_CMD) $$SIGN_KEY_FILE $(BIN_DIR)/$(APP_NAME)-linux-amd64); \
	LIN_ARM_SIG=$$($(SIGN_CMD) $$SIGN_KEY_FILE $(BIN_DIR)/$(APP_NAME)-linux-arm64); \
	WIN_AMD_SIG=$$($(SIGN_CMD) $$SIGN_KEY_FILE $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe); \
	WIN_ARM_SIG=$$($(SIGN_CMD) $$SIGN_KEY_FILE $(BIN_DIR)/$(APP_NAME)-windows-arm64.exe); \
	MAC_AMD_SIG=$$($(SIGN_CMD) $$SIGN_KEY_FILE $(BIN_DIR)/$(APP_NAME)-darwin-amd64); \
	MAC_ARM_SIG=$$($(SIGN_CMD) $$SIGN_KEY_FILE $(BIN_DIR)/$(APP_NAME)-darwin-arm64); \
	UPDATED=$$(mktemp); \
	jq \
	--arg version "$(VERSION)" \
	--arg commit "$(COMMIT)" \
	--arg pubkey "$$(printf %s "$(PUBLIC_KEY)")" \
	--arg lin_amd_digest "$$LIN_AMD_DIGEST" \
	--arg lin_arm_digest "$$LIN_ARM_DIGEST" \
	--arg win_amd_digest "$$WIN_AMD_DIGEST" \
	--arg win_arm_digest "$$WIN_ARM_DIGEST" \
	--arg mac_amd_digest "$$MAC_AMD_DIGEST" \
	--arg mac_arm_digest "$$MAC_ARM_DIGEST" \
	--arg lin_amd_sig "$$LIN_AMD_SIG" \
	--arg lin_arm_sig "$$LIN_ARM_SIG" \
	--arg win_amd_sig "$$WIN_AMD_SIG" \
	--arg win_arm_sig "$$WIN_ARM_SIG" \
	--arg mac_amd_sig "$$MAC_AMD_SIG" \
	--arg mac_arm_sig "$$MAC_ARM_SIG" \
	-f scripts/merge_manifest.jq $(MANIFEST) > $$UPDATED; \
	mv $$UPDATED $(MANIFEST); \
	$(SIGN_CMD) $$SIGN_KEY_FILE $(MANIFEST) > $(MANIFEST).sig.base64; \
	rm -f $$SIGN_KEY_FILE; \
	echo "Manifest and signature updated: $(MANIFEST)"

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
