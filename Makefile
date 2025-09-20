include .env
export

run:
	docker compose up --build

run-d:
	docker compose up -d --build

migrate-create:
	docker run --rm \
		--user $(shell id -u):$(shell id -g) \
		-v $(PWD)/db/migrations:/migrations migrate/migrate \
		create -ext sql -dir /migrations -seq $(NAME)

migrate-up:
	docker run --rm \
		--network=modak-challenge_modak-challenge-net \
		-v $(PWD)/db/migrations:/migrations migrate/migrate \
		-path=/migrations \
		-database ${DB_URL} up

migrate-down:
	docker run --rm \
		--network=modak-challenge_modak-challenge-net \
		-v $(PWD)/db/migrations:/migrations migrate/migrate \
		-path=/migrations \
		-database ${DB_URL} down -all

migrate-force:
	docker run --rm \
		--network=modak-challenge_modak-challenge-net \
		-v $(PWD)/db/migrations:/migrations migrate/migrate \
		-path=/migrations \
		-database ${DB_URL} force $(VERSION)

schema-dump:
	docker run --rm \
		--network=modak-challenge_modak-challenge-net \
		-e PGPASSWORD=dev \
		${DB_IMAGE} \
		pg_dump ${DB_URL} \
		--schema-only --no-owner --no-privileges \
		> db/schema.sql

sqlc-generate:
	docker run --rm -v $(PWD):/src -w /src sqlc/sqlc:1.29.0 generate -f db/sqlc.yml

migrate-sync: migrate-up schema-dump sqlc-generate
	@echo "✅ Migration and SQLC sync completed successfully!"
