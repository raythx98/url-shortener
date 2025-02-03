createMigration:
	migrate create -ext sql -dir migrations/url_mappings -seq init_table

migrateUp:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	migrate -database ${POSTGRESQL_URL} -path migrations/url_mappings up

sqlc:
	sqlc generate --file url_mappings.yaml

swagger:
	swag init --parseDependency -o ./docs -d ./cmd/api,./controller

run:
	go run ./cmd/api/main.go