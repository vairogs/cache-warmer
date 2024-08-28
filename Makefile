## —— Cache Warmer —————————————————————————————————————————
help: ## Outputs this help screen
	@grep -E '(^[a-zA-Z0-9_-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m##/[33m/'

## —— Project ———————————————————————————————————————————————————————————————
run: ## Run the main go file on a Symfony project
	go run vcw.go $(path)

build: ## Build the vcw executable for the current platform
	go build -o bin/vcw vcw.go
	shasum -a 256 bin/vcw

clean: ## Clean all executable
	rm -f bin/vcw bin/vcw.exe

deps: clean ## Clean deps
	go get -d -v ./...

## —— Tests ✅ —————————————————————————————————————————————————————————————————
test: ## Run all tests
	go test -count=1 -v ./...

## —— Coding standards ✨ ——————————————————————————————————————————————————————
lint: ## Run gofmt simplify and lint
	gofmt -s -l -w .
