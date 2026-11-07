APP     ?= app
BIN_DIR ?= bin
GO      ?= go

GOFLAGS_BUILD := -trimpath -ldflags="-s -w"

.PHONY: all build test clean

all: test build

## build: compile a static binary into $(BIN_DIR)/$(APP)
build:
	CGO_ENABLED=0 $(GO) build $(GOFLAGS_BUILD) -o $(BIN_DIR)/$(APP) .

## test: run unit tests
test:
	$(GO) test -v ./...

clean:
	rm -rf $(BIN_DIR)
