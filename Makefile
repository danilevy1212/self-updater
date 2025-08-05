APP_NAME := api
CMD_DIR := ./cmd/$(APP_NAME)
BIN_DIR := bin

VERSION := $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse HEAD)

.PHONY: all build clean test format

all: format test build

build:
	go build -o $(BIN_DIR)/$(APP_NAME) \
		-ldflags "\
			-X 'main.Version=$(VERSION)' \
			-X 'main.Commit=$(COMMIT)'" \
		$(CMD_DIR)

test:
	go test -v ./... \
		-ldflags "\
			-X 'main.Version=$(VERSION)' \
			-X 'main.Commit=$(COMMIT)'"

clean:
	rm -rf $(BIN_DIR)/$(APP_NAME)

fmt:
	go fmt ./...
