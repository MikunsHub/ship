.PHONY: build install clean version help

VERSION := v0.1.0
BINARY_NAME := ship
INSTALL_PATH := $(shell go env GOPATH)/bin

help:
	@echo "Ship CLI - Build targets"
	@echo "  make build      - Build binary to ./bin/ship"
	@echo "  make install    - Install binary locally"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make version    - Show current version"

build:
	go build -ldflags "-X main.Version=$(VERSION)" -o bin/$(BINARY_NAME) .

install: build
	cp bin/$(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)

clean:
	rm -rf bin/

version:
	@echo "Ship CLI version $(VERSION)"
