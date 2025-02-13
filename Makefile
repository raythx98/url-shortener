createMigration:
	migrate create -ext sql -dir migrations/url_mappings -seq init_table

migrateUp:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	export POSTGRESQL_URL='postgres://postgres:password@localhost:5432/url_shortener?sslmode=disable'
	migrate -database ${POSTGRESQL_URL} -path migrations up

sqlc:
	cd sqlc
	sqlc generate --file sqlc.yaml
	cd ..

swagger:
	swag init --parseDependency -o ./docs -d ./cmd/api,./controller

run:
	go run ./cmd/api/main.go

#migrate -database postgres://postgres:password@localhost:5432/url_shortener?ssl_mode=disable -path ./migrations up