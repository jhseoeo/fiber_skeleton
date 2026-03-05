# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
make run       # go run ./cmd/main.go
make build     # go build -o bin/server ./cmd/main.go
make test      # go test ./...
make lint      # golangci-lint run
make swagger   # swag init -g cmd/main.go (regenerate docs/)
make tidy      # go mod tidy
```

Run a single test by name:
```bash
go test ./src/handler/... -run TestGetExample_Success
```

Live reload (requires `air`):
```bash
air
```

Copy `.env.example` to `.env` before running. Config is loaded from environment variables (PORT, ENV, LOG_LEVEL, REQUEST_TIMEOUT, JWT_SECRET).
`ENV=production` 시 `JWT_SECRET`이 기본값이면 시작 시 panic합니다.

## Architecture

The project follows a strict layered architecture: **handler → service → repository**, with each layer communicating only through port interfaces defined in the layer below.

```
cmd/main.go                  # wires everything together
src/
  config/                    # env-based config struct
  handler/                   # HTTP layer; depends on serviceport interfaces
  service/                   # business logic; depends on repositoryport interfaces
    serviceport/             # port interfaces consumed by handlers
  repository/                # data layer (currently in-memory placeholder)
    port/                    # port interfaces consumed by services
    error/                   # sentinel errors (ErrNotFound, ErrAlreadyExists)
  middleware/                # error_handler, logger, recover, auth, security, timeout
  model/                     # domain models
  dto/
    errorcode/               # custom error codes
    req/                     # request DTOs (with go-playground/validator tags)
    resp/                    # response DTOs (CommonResp, PaginatedResp)
  pkg/
    typeerr/                 # ErrorResp type + SentinelError
    validate/                # go-playground/validator wrapper
    log/                     # logrus setup
docs/                        # generated swagger stub (do not edit manually)
```

## Infrastructure Files

- `Makefile` — common dev commands
- `Dockerfile` — multi-stage build (builder: `golang:1.26-alpine`, runtime: `alpine:3.21`)
- `docker-compose.yml` — single-service compose referencing `.env`
- `.air.toml` — live reload config; builds to `tmp/main`

## Key Patterns

### Error handling flow
Handlers return `typeerr.NewErrorResp(err, errorcode.ErrXxx, "message")`, which is caught by `middleware.NewErrorHandler()` (registered as `fiber.Config.ErrorHandler`). The error handler derives the HTTP status from `ErrorCode.HTTPStatus()` which is `int(code)/100` — so error code `40400` → HTTP 404.

### Adding new error codes
Add constants to `src/dto/errorcode/errorcode.go`. The numeric value encodes the HTTP status: `XXYYZZ` where `XXY` = HTTP status code (first 3 digits), e.g., `40401` → HTTP 404.

### Adding new repository sentinel errors
Add to `src/repository/error/error.go` using `typeerr.NewSentinelError("message")`. Check with `errors.Is(err, repositoryerror.ErrNotFound)` in handlers.

### Interface compliance check
Services and repositories declare compile-time interface compliance:
```go
var _ serviceport.ExampleServicePort = (*ExampleService)(nil)
```

### Testing pattern
- **Mock types** live in `src/testutil/` — `MockExampleService` and `MockExampleRepository` with function fields per method.
- **Service tests** (`src/service/*_test.go`): inject `MockExampleRepository`, no HTTP.
- **Handler unit tests** (`src/handler/example_test.go`): inject `MockExampleService` via `newTestApp()`.
- **Handler integration tests** (`src/handler/example_integration_test.go`): wire real `ExampleRepository` + `ExampleService` through `newIntegrationApp()`.
- All `app.Test` calls use `fiber.TestConfig{Timeout: 5 * time.Second}`.

### Validation errors
`validate.Struct()` returns `validate.FieldErrors` (a `[]FieldError{Field, Message}`) on failure.
In handlers, use `typeerr.NewErrorRespWithData(err, code, message, err)` to include the structured list in the response `data` field.

### JWT auth
`middleware.NewAuthMiddleware([]byte(secret))` validates `Authorization: Bearer <token>`. Access claims via `middleware.GetClaims(c)`. Routes are not protected by default — wrap route groups with the middleware when needed (see the TODO comment in `cmd/main.go`).

### Middleware stack order (cmd/main.go)
`Recoverer → Metrics → CORS → Security (helmet) → Timeout → RequestID → Logger`
- **Metrics** (`NewMetrics(app)`) also registers `GET /metrics` (Prometheus) on the app.
- **CORS** (`NewCORS()`) defaults to allow-all; pass `cors.Config{AllowOrigins: ...}` in production.
- **RateLimiter** (`NewRateLimiter(limiter.Config{...})`) is not wired by default — apply to specific route groups.

### Health endpoints
- `GET /health/live` — liveness probe (always 200 if process is running)
- `GET /health/ready` — readiness probe (add dependency checks here; return 503 when unavailable)

### Swagger
Annotations live on unexported handler methods. Run `make swagger` to regenerate `docs/docs.go`. The UI is at `GET /swagger`, spec at `GET /swagger/doc.json`.

### pprof
Only enabled when `ENV=development`. Available at `/debug/pprof`.
