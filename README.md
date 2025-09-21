# Modak Challenge - Notification API

Notification API built with Go, Gin, and PostgreSQL. It supports per-user, per-type rate limiting and persists notification events. The project follows a clean, hexagonal architecture (ports and adapters), is fully containerized with Docker, and ships with Swagger docs, SQLC for type-safe queries, and Make targets for common tasks.

## Key Features

- **Per-type rate limiting** per user, enforced in the domain use case
  - `status`: 2 per minute
  - `news`: 1 per 24 hours
  - `marketing`: 3 per hour
- **Persistent storage** of notifications in PostgreSQL with efficient index for time-window queries
- **HTTP API** using Gin with health check and Swagger UI
- **Hexagonal architecture** separating use case, ports, and adapters
- **SQLC** generated database access for type-safe queries
- **Docker Compose** for local development, plus optional pgAdmin
- **CI** workflow to validate SQLC generation and build/push Docker images

## Architecture

Code is organized following ports-and-adapters:

- `internal/domain/usecase/`: business rules and rate limiting
- `internal/ports/`: interfaces used by the use case
- `internal/adapters/db/`: Postgres repository implementation backed by SQLC
- `internal/adapters/gateway/`: outbound notification gateway (a sample gateway that prints to console)
- `internal/adapters/http/`: HTTP server, routing, handlers, and DTOs
- `internal/config/`: app config and domain errors

## Tech Stack

- **Language:** Go 1.23
- **HTTP:** Gin
- **Database:** PostgreSQL (with SQLC-generated queries)
- **Docs:** Swagger (swaggo)
- **Containerization:** Docker, Docker Compose
- **CI:** GitHub Actions

## Project Structure

```
cmd/server/main.go                  # App entrypoint (wire adapters and use case)
internal/
  adapters/
    http/                           # Router, routes v1, handlers and DTOs
    db/                             # SQLC repo implementation
    gateway/                        # Fake notification gateway (console)
  config/                           # App config and domain errors
  domain/
    entity/                         # Entities and rate-limit configuration
    usecase/                        # Core business logic (rate limiting)
  ports/                            # Interfaces (repository, gateway)
db/
  migrations/                       # SQL migrations
  queries/                          # SQLC input queries
  schema.sql                        # Schema dump (source of truth for SQLC)
```

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Make (optional but recommended)
- Go 1.23+ (only if running locally without Docker)

### Environment Variables

Copy `.env.example` to `.env` at the root of the project and set values. Typical local values:

```env
DB_USER=dev
DB_PASSWORD=dev
DB_NAME=ModakChallengeDB
DB_PORT=5432
DB_IMAGE=postgres:16.8-alpine3.20
DB_VOLUME=.data
DB_HOST=modak-challenge-db

# Connection string used by the API and migrations
DB_URL=postgresql://dev:dev@modak-challenge-db:5432/ModakChallengeDB?sslmode=disable

PGADMIN_IMAGE=dpage/pgadmin4:9.4
PGADMIN_EMAIL=dev@modakchallenge.com
PGADMIN_PASSWORD=dev
PGADMIN_PORT=5050
PGADMIN_HOST=modak-challenge-pgadmin

APP_PORT=8080
GO_VERSION=1.23-alpine
```

Notes:

- `DB_URL` must be valid for both the API container and migration commands.

### Quickstart (Docker Compose)

Start the full stack (API + Postgres + pgAdmin):

```bash
make run            # or: docker compose up --build
# or detached
make run-d          # or: docker compose up -d --build
```

Apply migrations, dump schema, and regenerate SQLC (recommended):

```bash
make migrate-up
make schema-dump
make sqlc-generate
# or in one shot:
make migrate-sync
```

Once running:

- Health check: http://localhost:8080/health
- Swagger UI: http://localhost:8080/swagger/index.html
- pgAdmin: http://localhost:5050 (use the credentials from `.env`)

### Local development without Docker (optional)

1) Start a local PostgreSQL and create the database from the migrations.

2) Export `DB_URL` and `APP_PORT`, then run:

```bash
go run ./cmd/server/main.go
```

## API

### Health Check

- Method: `GET /health`
- Response: `200 {"status":"ok"}`

### Send Notification

- Method: `POST /v1/notifications/send`
- Request body:

```json
{
  "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "type": "status",    
  "message": "Your order shipped"
}
```

- Success: `201 {"status":"sent"}`
- Errors:
  - `400 {"error":"invalid notification"}` for unsupported `type` or validation errors
  - `429 {"error":"rate limit exceeded"}` when the per-type limit is reached
  - `500` for unexpected server/database issues

Example request:

```bash
curl -X POST http://localhost:8080/v1/notifications/send \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
    "type": "status",
    "message": "Your order shipped"
  }'
```

### Rate Limits

- `status`: 2 notifications per 1 minute
- `news`: 1 notification per 24 hours
- `marketing`: 3 notifications per 1 hour

The check is performed in the use case by counting messages for the user and type since a computed window start and comparing it to the limit.

## Database

- Table: `notifications`
  - Columns: `id (uuid)`, `user_id (uuid)`, `type (text)`, `message (text)`, `created_at (timestamp)`
  - Index: `idx_notifications_user_type_time` on `(user_id, type, created_at)` to serve the time-window count efficiently
- SQLC:
  - Queries in `db/queries/notifications.sql`
  - Code generated to `internal/adapters/db/sqlc` using `db/sqlc.yml`

Useful Make targets:

```bash
make migrate-up        # Apply all migrations
make migrate-down      # Roll back all migrations (danger: all)
make migrate-force     # Force set migration version
make migrate-version   # Show current migration version
make schema-dump       # Dump current DB schema into db/schema.sql
make sqlc-generate     # Generate SQLC code from schema + queries
make migrate-sync      # migrate-up + schema-dump + sqlc-generate
```

## Development Tooling

- Swagger docs are served at `/swagger/*` and generated with:

```bash
make swagger-generate
```

This runs `swag init` inside a Go container and outputs docs into `./docs`. The server imports these docs in `cmd/server/main.go`.

## Testing

Run tests and open coverage report:

```bash
make test           # runs unit tests across modules
```

The coverage visualization step (`make test-cover`) filters out non-essential files and opens an HTML report. The HTML will be generated in a temporary folder, so you need follow the provided path to view it.

## CI/CD

GitHub Actions workflow (`.github/workflows/build.yml`):

- **sqlc job**: Generates SQLC code using Docker and fails if changes are not committed.
- **test job**: Runs tests and fails if any test not passing.
- **build job**: Builds and pushes Docker images to Docker Hub as `latest` and `SHA` tags.

To push images, set the `DOCKERHUB_TOKEN` repository secret. The username and image name are configured via workflow env vars (`DOCKERHUB_USERNAME=paulooo`, `IMAGE_NAME=modak-challenge`). Link to the dockerhub repository [here](https://hub.docker.com/repository/docker/paulooo/modak-challenge).

## Troubleshooting

- **Database connection failed**: Ensure `DB_URL` points to the Compose service (`modak-challenge-db`) and the database is up.
- **Migrations not applied**: Run `make migrate-up` after the database is running. Then `make schema-dump && make sqlc-generate`.
- **Migration version dirty**: If applyed a break change in migrations and it's dirty, run `make migrate-force` to force another migration version, then `make migrate-sync` to apply migrations and regenerate SQLC.
- **Swagger not found**: Run `make swagger-generate` to create the `docs/` folder.
- **Rate limit unexpected**: Verify the notification `type` is one of `status`, `news`, or `marketing` and that the time window matches your expectations.