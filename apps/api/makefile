.PHONY: help
help:
	@echo "api       - run http server"
	@echo "worker    - run outbox worker"
	@echo "cron      - run cron scheduler"
	@echo "build     - build all binaries"
	@echo "test      - run tests"
	@echo "lint      - go vet + staticcheck"
	@echo "migrate   - run migrations up"
	@echo "sqlc      - generate sqlc code"
	@echo "tidy      - go mod tidy"

.PHONY: api worker cron
api:
	go run ./cmd/api
worker:
	go run ./cmd/worker
cron:
	go run ./cmd/cron

.PHONY: build
build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/api ./cmd/api
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/worker ./cmd/worker
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/cron ./cmd/cron

.PHONY: test
test:
	go test -race -count=1 ./...

.PHONY: lint
lint:
	go vet ./...
	@command -v staticcheck >/dev/null 2>&1 && staticcheck ./... || echo "staticcheck not installed, skipping"

.PHONY: migrate
migrate:
	migrate -path ./migrations -database "$$MIGRATIONS_URL" up

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: tidy
tidy:
	go mod tidy