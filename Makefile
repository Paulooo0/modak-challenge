include .env
export

run:
	docker compose up --build

run-d:
	docker compose up -d --build

test: test-run test-cover

test-run:
	go test ./... -coverpkg=./... -coverprofile=coverage.out

test-cover:
	grep -v "cmd/server" coverage.out | \
	grep -v "internal/adapters/db/sqlc" | \
	grep -v "internal/adapters/http/server.go" | \
	grep -v "internal/adapters/http/v1/routes.go" | \
	grep -v "_router.go" | \
	grep -v "internal/config/config.go" > coverage.filtered.out
	go tool cover -html=coverage.filtered.out

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

migrate-version:
	docker run --rm \
		--network=modak-challenge_modak-challenge-net \
		-v $(PWD)/db/migrations:/migrations migrate/migrate \
		-path=/migrations \
		-database ${DB_URL} version

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

swagger-generate:
	docker run --rm -v $(PWD):/src -w /src -v swag-go-cache:/go golang:1.23-alpine sh -c 'apk add --no-cache git >/dev/null && go install github.com/swaggo/swag/cmd/swag@v1.16.4 && /go/bin/swag init -g cmd/server/main.go -o docs'