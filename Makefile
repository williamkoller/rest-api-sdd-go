.PHONY: run build migrate docker-up docker-down tidy fmt lint vulnerability tools hooks

COMPOSE := docker-compose -f docker/docker-compose.yml

run:
	$(COMPOSE) up postgres redis -d
	go run ./cmd/api/...

build:
	go build -o bin/api ./cmd/api/...

migrate:
	$(COMPOSE) run --rm migrate

docker-up:
	$(COMPOSE) up --build

docker-down:
	$(COMPOSE) down

tidy:
	go mod tidy

fmt:
	gofmt -s -w .

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run

vulnerability:
	govulncheck ./...

tools:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	@echo "Install golangci-lint: brew install golangci-lint"
	@echo "Install lefthook:      brew install lefthook"

hooks:
	lefthook install
