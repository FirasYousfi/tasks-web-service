BINARY=engine

.PHONY: run
run: ## run the tasks-web-service app
	go run ./cmd/server/main.go

test: ## run the unit tests
	go test ./...

swagger_generate: ## generate the swagger documentation files
	go install github.com/swaggo/swag/cmd/swag@latest && \
	swag init --parseDependency -d cmd/server

swagger_ui: ## generate the swagger documentation files
	make swagger_generate && \
	go run ./cmd/swagger/main.go

run_postgres: ## spin up a postgresql database with docker
	docker compose up

stop_postgres: ## stop the postgresql database
	docker compose down

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)