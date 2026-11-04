# Go API Skeleton

Small REST API base built with Go, Echo, sqlx, and Viper. A fresh clone starts without a config file or database.

## What is included

- Graceful shutdown for `SIGINT` and `SIGTERM`
- HTTP read, write, idle, and header timeouts
- Liveness and dependency-aware readiness probes
- JSON request logs with request IDs and latency
- Panic recovery, a 2 MiB body limit, CORS, CSP, and security headers
- Optional MySQL connection verified during startup
- Static landing page served by the API
- Unit tests, CI, and a non-root Docker image

## Requirements

- Go 1.27+
- Docker with Compose, only if running the container workflow
- MySQL, only if database support is enabled

## Quick start

```bash
git clone https://github.com/Ikuradachi/skeleton_go.git
cd skeleton_go
go run ./cmd/api
```

Verify the process:

```bash
curl http://localhost:38600/health/ready
```

Expected response without a configured database:

```json
{
  "status": "ok",
  "checks": {
    "database": "disabled"
  }
}
```

Open <http://localhost:38600> for the landing page.

## Routes

| Method | Path | Purpose | Success |
| --- | --- | --- | --- |
| `GET` | `/api` | API identity and status | `200` |
| `GET` | `/health/live` | Confirms the process can serve HTTP | `200` |
| `GET` | `/health/ready` | Checks configured dependencies | `200` or `503` |

Unknown routes return a stable JSON error with the request ID:

```json
{
  "error": {
    "code": "not_found",
    "message": "resource not found",
    "request_id": "..."
  }
}
```

## Configuration

Configuration is applied in this order:

1. Built-in defaults
2. Optional `pkg/config/config.yaml`
3. `SKELETON_*` environment variables

Copy the example only when values need to change:

```bash
cp pkg/config/config.yaml.example pkg/config/config.yaml
```

`pkg/config/config.yaml` is ignored by Git because it may contain credentials. Setting `SKELETON_CONFIG_FILE` makes that explicit file required.

| Environment variable | Default | Meaning |
| --- | --- | --- |
| `SKELETON_CONFIG_FILE` | unset | Explicit YAML config path |
| `SKELETON_SERVER_HOST` | `0.0.0.0` | Listen address |
| `SKELETON_SERVER_PORT` | `38600` | Listen port |
| `SKELETON_DB_DRIVER` | empty | SQL driver; this repo includes `mysql` |
| `SKELETON_DB_DSN` | empty | Database connection string |
| `SKELETON_APP_PATH` | `app` | Static asset directory |
| `SKELETON_APP_FILE` | `index.html` | Root document |

### Enable MySQL

Set both database values. Setting only one fails validation before the HTTP server starts.

```yaml
db:
  driver: "mysql"
  dsn: "root:password@tcp(localhost:3306)/skeleton?parseTime=true"
```

Startup pings the database with a five-second timeout. `/health/ready` continues checking it after startup and returns `503` when it becomes unavailable.

Each instance opens at most 10 database connections. Requests wait for an available connection when the pool is full, bounded by their context deadline.

## Project layout

```text
.
├── cmd/api/main.go              # composition root and process lifecycle
├── internal/api/                # router, middleware, HTTP errors
├── internal/features/health/    # example vertical feature slice
├── pkg/config/                  # config loading and database connection
├── pkg/logger/                  # slog JSON logger
├── app/                         # static landing page
├── Dockerfile
└── docker-compose.yaml
```

Add a feature under `internal/features/<name>`. Keep its interfaces close to the consumer, register handlers in `internal/api/router.go`, and test the observable HTTP behavior. Add service or repository layers when the feature has logic that needs them.

## Development checks

```bash
go test ./...
go vet ./...
go build ./cmd/api
```

For live reload with [Air](https://github.com/air-verse/air), run `air`; `.air.toml` is already configured.

CORS currently permits anonymous `GET`, `HEAD`, and `OPTIONS` requests from any origin. Update `internal/api/middleware.go` before adding browser-based write endpoints or credentials.

## Docker

```bash
docker compose up --build
```

Set a different host port when needed:

```bash
APP_PORT=8080 docker compose up --build
```

The container runs as an unprivileged user with a read-only filesystem, dropped Linux capabilities, and an HTTP health check. To use a config file, create it first and add this volume to the `api` service:

```yaml
volumes:
  - ./pkg/config/config.yaml:/app/pkg/config/config.yaml:ro
```

## License

[MIT](LICENSE)
