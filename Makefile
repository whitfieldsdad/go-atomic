BUILD_DIRECTORY=bin
FILENAME_PREFIX=go-atomic-red-team
LDFLAGS = "-s -w"

ATOMIC_RED_TEAM_PATH=~/src/atomic-red-team

all: help

help:
	@echo "build - Build all binaries"
	@echo "update - Update data files"
	@echo "docs - Generate documentation"
	@echo "[test|tests] - Run tests"
	@echo "windows - Build Windows binaries"
	@echo "linux - Build Linux binaries"
	@echo "darwin - Build Darwin binaries"
	@echo "386 - Build 386 binaries"
	@echo "amd64 - Build amd64 binaries"
	@echo "arm64 - Build arm64 binaries"

bin: build

docs:
	make -C docs

build: windows linux darwin
	chmod +x bin/*
	du -sh bin/*

update:
	go run main.go task-templates -w data/content generate $(ATOMIC_RED_TEAM_PATH)/atomics

test:
	go test -v ./...

tests: test

windows: windows_amd64 windows_arm64
linux: linux_386 linux_amd64 linux_arm64
darwin: darwin_amd64 darwin_arm64

386: linux_386
amd64: windows_amd64 linux_amd64 darwin_amd64
arm64: windows_arm64 linux_arm64 darwin_arm64

windows_amd64:
	GOOS=windows GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o $(BUILD_DIRECTORY)/$(FILENAME_PREFIX)-windows-amd64.exe .

windows_arm64:
	GOOS=windows GOARCH=arm64 go build -ldflags=$(LDFLAGS) -o $(BUILD_DIRECTORY)/$(FILENAME_PREFIX)-windows-arm64.exe .

linux_386:
	GOOS=linux GOARCH=386 go build -ldflags=$(LDFLAGS) -o $(BUILD_DIRECTORY)/$(FILENAME_PREFIX)-linux-386 .

linux_amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o $(BUILD_DIRECTORY)/$(FILENAME_PREFIX)-linux-amd64 .

linux_arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags=$(LDFLAGS) -o $(BUILD_DIRECTORY)/$(FILENAME_PREFIX)-linux-arm64 .

darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags=$(LDFLAGS) -o $(BUILD_DIRECTORY)/$(FILENAME_PREFIX)-darwin-amd64 .

darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags=$(LDFLAGS) -o $(BUILD_DIRECTORY)/$(FILENAME_PREFIX)-darwin-arm64 .

.PHONY: docs build