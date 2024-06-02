#export POSTGRESQL_URL='postgres://postgres:password@localhost:5432/example?sslmode=disable'

createMigration:
	migrate create -ext sql -dir migrations/url_mappings -seq init_table

runMigration:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	migrate -database ${POSTGRESQL_URL} -path migrations/url_mappings up

sqlcGenerate:
	sqlc generate --file url_mappings.yaml

#go get "github.com/raythx98/gohelpme"