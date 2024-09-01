default: help

SHELL:= /usr/bin/bash -e

BINDIR := $(CURDIR)/bin
BINNAME ?= gobackup
GOLANGCI_LINT_VERSION:=1.52.0
# Download the linter executable file from 
# https://github.com/golangci/golangci-lint/releases/tag/v1.52.2
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
CURRENT_TIME = $(shell date -u '+%Y-%m-%d %I:%M:%S %p GMT')

LDFLAGS = -X 'main.GitCommit=${GIT_SHA}' \
				-X 'main.GitTag=${GIT_TAG}' \
				-X 'main.BuildTime=${CURRENT_TIME}'

EXTLDFLAGS = "-extldflags '-static' ${LDFLAGS}"

.PHONY: help
help: ## Display this help details.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: all
all: lint test build ## Run linters then unit tests and build the app.

.PHONY: install-linter
install-linter: ## Install golangci-lint tool.
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v${GOLANGCI_LINT_VERSION}

.PHONY: lint
lint: ## Clean package and fix linters warnings.
	## use `make install-linter` to install linters if missing
	## or download the executable file from below link
	## https://github.com/golangci/golangci-lint/releases/
	go mod tidy
	golangci-lint run -v --skip-dirs bin

.PHONY: clean-test
clean-test: ## Remove temporary files and cached tests results.
	go clean -testcache

.PHONY: test
test: clean-test ## Remove cache and Run unit tests only.
	go test -v ./... -count=1

.PHONY: test-cover
test-cover: clean-test ## Remove tests cache and runs unit tests in non-verbose mode with coverage.
	go test ./... -cover -count=1

.PHONY: coverage-console
coverage-console: ## Obtain codebase testing coverage and view stats in console.
	go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out

.PHONY: coverage-web
coverage-web: ## Obtain codebase testing coverage and view stats in browser.
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

.PHONY: coverage
coverage: coverage-console coverage-web ## View codebase testing coverage on console and browser.

.PHONY: clean-build
clean-build: ## Remove temporary and cached builds files.
	go clean -cache

.PHONY: build
build: clean-build ## Remove builds cache and build the executable.
	CGO_ENABLED=0 go build -o ${BINDIR}/${BINNAME} -a -ldflags ${EXTLDFLAGS} main.go
	
.PHONY: format
format: ## Format the codebase.
	gofumpt -l -w .