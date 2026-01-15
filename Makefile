.PHONY: build run test clean create_migration sqlc swagger run_local up down logs

# Default migration name if not provided (make create_migration name=my_migration)
name ?= new_migration

build:
	go build -o bin/server cmd/api/main.go

run: build
	./bin/server

test:
	go test -v ./...

clean:
	rm -rf bin/

create_migration:
	migrate create -ext sql -dir migrations -seq $(name)

sqlc:
	cd sqlc && sqlc generate --file sqlc.yaml

swagger:
	swag init --parseDependency --parseInternal --parseDepth 2 -o ./docs -d ./cmd/api,./controller,./dto

run_local: sqlc swagger
	set -a && . .envrc && set +a && go run cmd/api/main.go

# Docker Compose Helpers
up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f
