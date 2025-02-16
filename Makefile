create_migration:
	migrate create -ext sql -dir migrations/url_mappings -seq init_table

sqlc:
	cd sqlc
	sqlc generate --file sqlc.yaml
	cd ..

swagger:
	swag init --parseDependency -o ./docs -d ./cmd/api,./controller

allow_direnv:
	direnv allow .

build:
	docker build -t url_shortener .

volume:
	docker volume create local-postgres

network:
	docker network create my-network || true

db:
	docker run -d --rm --name ${APP_URLSHORTENER_DBHOST} \
		--net my-network -p ${APP_URLSHORTENER_DBPORT}:${APP_URLSHORTENER_DBPORT} \
		-e POSTGRES_PASSWORD=${APP_URLSHORTENER_DBPASSWORD} \
		-v local-postgres:/var/lib/postgresql/data \
		postgres:latest || true

format_env:
	sed 's/^export //' .envrc > .env

migrate_up:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	migrate -database 'postgres://${APP_URLSHORTENER_DBUSERNAME}:${APP_URLSHORTENER_DBPASSWORD}@localhost:${APP_URLSHORTENER_DBPORT}/${APP_URLSHORTENER_DBDEFAULTNAME}?sslmode=disable' -path migrations up

run: allow_direnv build volume network db format_env migrate_up
	docker run -d --rm --name url-shortener-app \
		--net my-network -p ${APP_URLSHORTENER_SERVERPORT}:${APP_URLSHORTENER_SERVERPORT} \
		--env-file .env url-shortener

stop:
	docker stop url-shortener-app url-shortener-db