.PHONY: mod clean build start d-start d-stop dynamodb-admin help

# Show usage if no option
.DEFAULT_GOAL := help

mod: ## Get dependencies
	go get -u ./...

clean: ## Remove binary
	rm -rf ./hello-world/hello-world
	
build: ## Build binary
	GOOS=linux GOARCH=amd64 go build -o hello-world/hello-world ./hello-world

start: build ## Start Lambda on localhost
	sam local start-api # or start-lambda

d-start: ## Boot DynamoDB local
	docker-compose up -d

d-stop: ## Treminate DynamoDB local
	docker-compose down

pre-admin:
	@if [ -z `which dynamodb-admin2 2> /dev/null` ]; then \
		echo "Need to install dynamodb-admin, execute \"npm install dynamodb-admin -g\"";\
		exit 1;\
	fi

dynamodb-admin: pre-admin ## Start DaynamoDB GUI
	DYNAMO_ENDPOINT=http://localhost:18000 dynamodb-admin

help: ## Show options
	@grep -E '^[a-zA-Z_\.-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

