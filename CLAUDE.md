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

Copy `.env.example` to `.env` before running. Config is loaded from environment variables:
`PORT`, `ENV`, `LOG_LEVEL`, `REQUEST_TIMEOUT`, `JWT_SECRET`, `CORS_ALLOW_ORIGINS`.

### Production config constraints

- `ENV=production` 시 `JWT_SECRET`이 기본값이면 시작 시 panic.
- `JWT_SECRET`은 production에서 최소 32바이트 필요 (HMAC 보안).
- `CORS_ALLOW_ORIGINS`가 비어 있으면 전체 오리진 허용 (production에서 경고 로그 출력).
- `REQUEST_TIMEOUT`이 유효하지 않은 값이면 경고 로그 후 기본값(30s) 사용.

## Architecture

The project follows a strict layered architecture: **handler → service → repository**, with each layer communicating only through port interfaces defined in the layer below.

```text
cmd/main.go                  # wires everything together
src/
  config/                    # env-based config struct with production validation
  handler/                   # HTTP layer; depends on serviceport interfaces
  service/                   # business logic; depends on repositoryport interfaces
    serviceport/             # port interfaces consumed by handlers
  repository/                # data layer (currently in-memory with sync.RWMutex)
    port/                    # port interfaces consumed by services
    error/                   # sentinel errors (ErrNotFound, ErrAlreadyExists)
  middleware/                # error_handler, logger, recover, auth, security, timeout, metrics, cors, ratelimit
  model/                     # domain models
  dto/
    errorcode/               # custom error codes (ErrBadRequest, ErrNotFound, ErrTooManyRequests, etc.)
    req/                     # request DTOs (with go-playground/validator tags, max constraints)
    resp/                    # response DTOs (CommonResp, PaginatedResp)
  pkg/
    typeerr/                 # ErrorResp type + SentinelError
    validate/                # go-playground/validator wrapper returning FieldErrors
    log/                     # logrus setup with request-scoped fields
  testutil/                  # shared mock types (MockExampleService, MockExampleRepository)
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

### Binding helpers

Handlers use `bindJSON(c, dst)` and `bindQuery(c, dst)` helpers that combine Fiber binding + `validate.Struct()` in one call. `bindJSON` failures use `ErrInvalidBody`; `bindQuery` failures use `ErrBadRequest`.

### Adding new error codes

Add constants to `src/dto/errorcode/errorcode.go`. The numeric value encodes the HTTP status: `XXYYZZ` where `XXY` = HTTP status code (first 3 digits), e.g., `40401` → HTTP 404. Current codes: `ErrBadRequest(40000)`, `ErrNotFound(40400)`, `ErrRequestTimeout(40800)`, `ErrConflict(40900)`, `ErrTooManyRequests(42900)`, `ErrInternalServer(50000)`.

### Adding new repository sentinel errors

Add to `src/repository/error/error.go` using `typeerr.NewSentinelError("message")`. Check with `errors.Is(err, repositoryerror.ErrNotFound)` in handlers.

### Repository concurrency safety

The in-memory repository uses `sync.RWMutex` and always stores/returns **value copies** of `model.Example` to prevent shared-pointer races after lock release. Follow this pattern when adding new repositories.

### Interface compliance check

Services and repositories declare compile-time interface compliance:

```go
var _ serviceport.ExampleServicePort = (*ExampleService)(nil)
```

### Testing pattern

- **Mock types** live in `src/testutil/` — `MockExampleService` and `MockExampleRepository` with function fields per method. Nil function fields panic with a descriptive message.
- **Service tests** (`src/service/*_test.go`): inject `MockExampleRepository`, no HTTP.
- **Handler unit tests** (`src/handler/example_test.go`): inject `MockExampleService` via `newTestApp()`. Covers all endpoints including List.
- **Handler integration tests** (`src/handler/example_integration_test.go`): wire real `ExampleRepository` + `ExampleService` through `newIntegrationApp()`. Each test creates an isolated app instance.
- All `app.Test` calls use `fiber.TestConfig{Timeout: 5 * time.Second}`.

### Validation

- `validate.Struct()` returns `validate.FieldErrors` (a `[]FieldError{Field, Message}`) on failure.
- In handlers, use `typeerr.NewErrorRespWithData(err, code, message, err)` to include the structured list in the response `data` field.
- Request DTOs enforce upper bounds (e.g., `Page max=10000`, `Limit max=100`, `Content max=10000`).

### JWT auth

`middleware.NewAuthMiddleware([]byte(secret))` validates `Authorization: Bearer <token>`. Access claims via `middleware.GetClaims(c)`. Routes are not protected by default — wrap route groups with the middleware when needed (see the TODO comment in `cmd/main.go`).

### Middleware stack order (cmd/main.go)

`Recoverer → Metrics → CORS → Security (helmet) → Timeout → RequestID → Logger`

- **Metrics** (`NewMetrics(app)`) registers Prometheus counter/histogram and `GET /metrics`. Uses `c.Route().Path` (with nil guard) to avoid label cardinality explosion.
- **CORS** (`NewCORS()`) defaults to allow-all with `MaxAge=3600`. `CORS_ALLOW_ORIGINS` env var accepts comma-separated origins. Empty entries from trailing/consecutive commas are filtered out.
- **RateLimiter** (`NewRateLimiter(limiter.Config{...})`) is not wired by default — apply to specific route groups.
- **Timeout** returns 408 when `context.DeadlineExceeded`, regardless of handler return value.

### Health endpoints

- `GET /health/live` — liveness probe (always 200 if process is running)
- `GET /health/ready` — readiness probe (add dependency checks here; return 503 when unavailable)

### Catch-all 404

`handler.NotFound` is registered as the last middleware and returns a JSON `CommonResp` with `errorcode.ErrNotFound` for any unregistered route.

### Swagger

Annotations live on unexported handler methods. Run `make swagger` to regenerate `docs/docs.go`. The UI is at `GET /swagger`, spec at `GET /swagger/doc.json`.

### pprof

Only enabled when `ENV=development`. Available at `/debug/pprof`.
