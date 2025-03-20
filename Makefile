create_migration:
	migrate create -ext sql -dir migrations/url_mappings -seq init_table

sqlc:
	cd sqlc
	sqlc generate --file sqlc.yaml
	cd ..

swagger:
	swag init --parseDependency -o ./docs -d ./cmd/api,./controller

allow_direnv:
	direnv allow . || true

build:
	docker build -t url-shortener .

volume:
	docker volume create local-postgres

network:
	docker network create my-network || true

db:
	docker run -d --rm --name ${APP_URLSHORTENER_DBHOST} \
		--net my-network -p ${APP_URLSHORTENER_DBPORT}:${APP_URLSHORTENER_DBPORT} \
		-e POSTGRES_PASSWORD=${APP_URLSHORTENER_DBPASSWORD} \
		-v local-postgres:/var/lib/postgresql/data \
		postgres:latest && sleep 5 || true


format_env:
	find .envrc && sed 's/^export //' .envrc > .env || true

migrate_up:
	migrate -database 'postgres://${APP_URLSHORTENER_DBUSERNAME}:${APP_URLSHORTENER_DBPASSWORD}@localhost:${APP_URLSHORTENER_DBPORT}/${APP_URLSHORTENER_DBDEFAULTNAME}?sslmode=disable' -path migrations up

run: allow_direnv swagger build volume network stop db format_env migrate_up
	docker run -d --rm --name url-shortener-app \
		--net my-network -p ${APP_URLSHORTENER_SERVERPORT}:${APP_URLSHORTENER_SERVERPORT} \
		--env-file .env url-shortener

stop_app:
	docker stop url-shortener-app || true

stop_db:
	docker stop url-shortener-db || true
	sleep 5

stop: stop_app stop_db