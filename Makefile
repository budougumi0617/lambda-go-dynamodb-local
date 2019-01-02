.PHONY: deps clean build start help

# Show usage if no option
.DEFAULT_GOAL := help

deps: ## Get dependencies
	go get -u ./...

clean: ## Remove binary
	rm -rf ./hello-world/hello-world
	
build: ## Build binary
	GOOS=linux GOARCH=amd64 go build -o hello-world/hello-world ./hello-world

start: build ## Start Lambda on localhost
	sam local start-api # or start-lambda

help: ## Show options
	@grep -E '^[a-zA-Z_-{\.}]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'

