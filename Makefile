# Docker related variables.
SERVICE_NAME = api
DB_SERVICE = postgres
DOCKER_COMPOSE = docker-compose

# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin

.PHONY: build up down clean test lint logs

all: build up

build:
	@echo "Building Docker images..."
	@$(DOCKER_COMPOSE) build $(SERVICE_NAME)

up:
	@echo "Starting services..."
	@$(DOCKER_COMPOSE) up -d

up_database:
	@echo "Starting services..."
	@$(DOCKER_COMPOSE) up -d $(DB_SERVICE)

down:
	@echo "Stopping and removing containers..."
	@$(DOCKER_COMPOSE) down -v

clean:
	@echo "Cleaning API resources..."
	@$(DOCKER_COMPOSE) stop $(SERVICE_NAME) 2>/dev/null || true
	@$(DOCKER_COMPOSE) rm -v -f $(SERVICE_NAME) 2>/dev/null || true
	@docker rmi -f $$(docker images -q mini-url-api) 2>/dev/null || true
	@echo "API cleanup complete - Postgres remains untouched"


prune: down
	@echo "Cleaning up resources..."
	@docker volume prune -f
	@docker system prune -f
logs:
	@$(DOCKER_COMPOSE) logs -f $(SERVICE_NAME)

db-shell:
	@$(DOCKER_COMPOSE) exec $(DB_SERVICE) psql -U postgres -d mini_url


# GO
go-compile: go-check-dependencies go-build-service
## Format
code-clean: go-format go-clean-dependencies

code-quality: go-static-checks go-vet go-format go-clean-dependencies

go-vet:
	@echo "  >  Running go vet..." # finds subtle issues where the code may not work as intended
	go vet ./internal/...

go-clean-dependencies:
	@echo "  > Running go mod tidy" # clean dependencies
	go mod tidy

go-format:
	@echo "  >  Running go formatter..." # Applies standard formatting (whitespace, indentation, etc)
	gofmt -w -s ./ 1>&2

go-static-checks: # Using static analysis, it finds bugs and performance issues, offers simplifications, and enforces style rules.
	@echo "  >  Run static checks..."
	@GOBIN=$(GOBIN) go install honnef.co/go/tools/cmd/staticcheck@master
	@$(GOBIN)/staticcheck -checks all ./...

go-test:
	@echo "  >  Run tests..."
	@GOBIN=$(GOBIN) go install github.com/onsi/ginkgo/v2/ginkgo
	@GOCOVERDIR=cover $(GOBIN)/ginkgo -r -race -cover --randomize-all --randomize-suites --trace --v ./

go-test-cov:
	@echo "  >  Run tests coverage..."
	go test ./internal/... -race -coverprofile=coverage.out -covermode=atomic ./ 1>&2
	go tool cover -html coverage.out

go-build-service:
	@echo "  >  Building service binary..."
	@GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/service $(GOBASE)/main.go

go-check-dependencies:
	@echo "  >  Checking if there are any missing dependencies..."
	go install ./...

go-install:
	@GOBIN=$(GOBIN) go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOBIN=$(GOBIN) go clean

go-run:
	@echo "  >  Install reflex"
	@GOBIN=$(GOBIN) go install github.com/cespare/reflex@latest
	@echo "  >  Running server..."
	@$(GOBIN)/reflex -R '^\.(go|idea)' -R '^(bin)' -r '.*\.go' -s -- /bin/sh -c 'go run $(GOBASE)/main.go'
