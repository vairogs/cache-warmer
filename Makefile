## —— Cache Warmer —————————————————————————————————————————
help: ## Outputs this help screen
	@grep -E '(^[a-zA-Z0-9_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

## —— Project ———————————————————————————————————————————————————————————————
run: ## Run the main go file on a Symfony project
	go run vcw.go $(path)

VERSION := $(shell git describe --tags --exact-match 2>/dev/null || git rev-parse HEAD)

build: ## Build the vcw executable for the current platform
	${MAKE} lint
	go build -ldflags="-X main.version=$(VERSION) -s -w" -o vcw
	strip vcw
	shasum -a 256 vcw

build-win: ## Build the vcw executable for the current platform
	go build -ldflags="-X main.version=$(VERSION) -s -w" -o vcw.exe
	shasum -a 256 vcw.exe

clean: ## Clean all executable
	rm -f vcw vcw.exe

deps: clean ## Clean deps
	go mod tidy
	go get -d -v ./...

update: ## Update dependecies
	go get -u ./...

## —— Tests ✅ —————————————————————————————————————————————————————————————————
test: ## Run all tests
	go test -count=1 -v ./...

## —— Coding standards ✨ ——————————————————————————————————————————————————————
lint: ## Run gofmt simplify and lint
	gofmt -s -l -w .
