.PHONY: help build run test clean docker-build docker-up docker-down security-check

help: ## Tampilkan help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build aplikasi
	go build -o api-user-crud-go .

run: ## Run aplikasi
	go run main.go

test: ## Run tests
	go test ./... -v -cover

test-coverage: ## Run tests dengan coverage report
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

clean: ## Clean build artifacts
	rm -f api-user-crud-go
	rm -f coverage.out
	rm -f test.db

docker-build: ## Build Docker image
	docker build -t user-crud-api .

docker-up: ## Start dengan docker-compose
	docker-compose up -d

docker-down: ## Stop docker-compose
	docker-compose down

docker-logs: ## View docker logs
	docker-compose logs -f

deps: ## Install dependencies
	go mod download
	go mod tidy

proto: ## Generate protobuf code
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user.proto

security-check: ## Run security checks
	@./scripts/security-check.sh
