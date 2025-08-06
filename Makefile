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

define sha256sum
$(shell sha256sum $(1) | cut -d ' ' -f1)
endef

.PHONY: all build-api clean test format release

all: format test build-api

build-api:
	go build -o $(BIN_DIR)/$(APP_NAME) \
		-ldflags "\
			-X 'main.Version=$(VERSION)' \
			-X 'main.Commit=$(COMMIT)'" \
		$(CMD_DIR)
	@echo "Build complete: $(BIN_DIR)/$(APP_NAME)"
	@echo "SHA256: $(call sha256sum,$(BIN_DIR)/$(APP_NAME))"

test:
	go test -v ./... \
		-ldflags "\
			-X 'main.Version=$(VERSION)' \
			-X 'main.Commit=$(COMMIT)'"

release: build-api
	@[ -n "$$$(SIGNING_KEY_ENV)" ] || (echo "FATAL: Environment variable $${SIGNING_KEY_ENV} is not set. Cannot sign release artifacts and manifest."; exit 1)

	@echo "Generating release manifest..."
	@DIGEST=$(call sha256sum,$(BIN_DIR)/$(APP_NAME)); \
	SIGN_KEY_FILE=$$(mktemp); \
	echo "$$$(SIGNING_KEY_ENV)" > $$SIGN_KEY_FILE; \
	SIG=$$($(SIGN_CMD) $$SIGN_KEY_FILE $(BIN_DIR)/$(APP_NAME)); \
	UPDATED=$$(mktemp); \
	jq \
		--arg version "$(VERSION)" \
		--arg commit "$(COMMIT)" \
		--arg filename "$(APP_NAME)-linux-amd64" \
		--arg digest "$$DIGEST" \
		--arg sig "$$SIG" \
		--arg pubkey "$$(printf %s "$(PUBLIC_KEY)")" \
		-f scripts/merge_manifest.jq $(MANIFEST) > $$UPDATED; \
	mv $$UPDATED $(MANIFEST); \
	$(SIGN_CMD) $$SIGN_KEY_FILE $(MANIFEST) > $(MANIFEST).sig.base64; \
	rm -f $$SIGN_KEY_FILE; \
	echo "Manifest and signature updated: $(MANIFEST)"

clean:
	rm -rf $(BIN_DIR)/$(APP_NAME)

format:
	go fmt ./...
