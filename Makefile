export CGO_ENABLED=0
SHELL    := /bin/bash

BINARY   := criminalsay
DIST     := dist
MAIN_PKG := .

.PHONY: all build test clean linux darwin-amd64 darwin-arm64 windows

all: linux darwin-amd64 darwin-arm64 windows

build:
	go build -o $(BINARY) $(MAIN_PKG)

test:
	go test ./...

clean:
	rm -rf $(DIST) $(BINARY)

linux: $(DIST)
	GOOS=linux GOARCH=amd64 go build -o $(DIST)/$(BINARY)-linux-amd64 $(MAIN_PKG)

darwin-amd64: $(DIST)
	GOOS=darwin GOARCH=amd64 go build -o $(DIST)/$(BINARY)-darwin-amd64 $(MAIN_PKG)

darwin-arm64: $(DIST)
	GOOS=darwin GOARCH=arm64 go build -o $(DIST)/$(BINARY)-darwin-arm64 $(MAIN_PKG)

windows: $(DIST)
	GOOS=windows GOARCH=amd64 go build -o $(DIST)/$(BINARY)-windows-amd64.exe $(MAIN_PKG)

$(DIST):
	mkdir -p $(DIST)
