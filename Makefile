GO ?= go

MAIN = main.go
BIN_DIR = ./bin
CLI_DIR = ${BIN_DIR}/aws-lambda-connector

# Detect the operating system
UNAME_S := $(shell uname -s)

ifeq ($(UNAME_S),Linux)
  OS_SUFFIX = -linux
endif

ifeq ($(UNAME_S),Darwin)
  OS_SUFFIX = -darwin
endif

ifeq ($(UNAME_S),CYGWIN*)
	OS_SUFFIX = .exe
endif

ifeq ($(UNAME_S),MINGW*)
	OS_SUFFIX = .exe
endif

ifeq ($(UNAME_S),Windows_NT)
	OS_SUFFIX = .exe
endif

EXEC = ${CLI_DIR}${OS_SUFFIX}

all:
	$(MAKE) dependencies
	$(MAKE) run

dependencies:
	$(GO) mod download
	$(GO) mod tidy

build-linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux $(GO) build -gcflags='all=-N -l' -o $(CLI_DIR)-linux $(MAIN)

build-darwin:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin $(GO) build -gcflags='all=-N -l' -o $(CLI_DIR)-darwin $(MAIN)

build-windows:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows $(GO) build -gcflags='all=-N -l' -o $(CLI_DIR).exe $(MAIN)

build: build${OS_SUFFIX}

clean:
	go clean
	rm -rf $(BIN_DIR)

.PHONY: all dependencies build run clean
