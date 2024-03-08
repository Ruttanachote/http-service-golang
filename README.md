# Stream HTTP Service Template (Version 2)

[Stream] Microservices REST API template using Go Fiber Framework + exception tracking by Sentry.io

## Extra feature includes

1. Object relational mapping (ORM)
2. In-memory Caching
3. Error exception handling
4. Observability (tracing)

## Documentation

- Fiber Framework (<https://docs.gofiber.io/>)
- GORM (<https://gorm.io/>)
- Go Redis (<https://redis.uptrace.dev/guide/go-redis.html>)
- Sentry (<https://docs.sentry.io/platforms/go/>)
- OpenTelemetry (<https://opentelemetry.io/docs/instrumentation/go/>)

## System requirements

- Golang 1.21
- PostgreSQL Database 15.x
- Redis 6

## Install dependencies

```bash
go install
```

## Update dependencies

```bash
go mod tidy
```

---

## Start development server

```bash
go run .
```

The server will listen by default at <http://localhost:8000>

---

## Build go executable file

```bash
go build -o main
```

## Environment Variables

```env
# App config
ENV="local"
HTTP_PORT=8000
GRPC_PORT=8001
APP_NAME="Stream - User Service"
SERVICE_NAME="user-service"

# Fiber config
FIBER_PREFORK=true
CORS_ALLOW_ORIGINS="*"
RATE_LIMIT=60

# Database config
DATABASE_HOST="localhost"
DATABASE_PORT=5432
DATABASE_NAME="postgres"
DATABASE_USER="postgres"
DATABASE_PASSWORD="mysecretpassword"
DATABASE_TIMEZONE="Asia/Bangkok"
DATABASE_MAX_IDLE_CONNS=2
DATABASE_MAX_OPEN_CONNS=3

# Redis config
REDIS_HOST="127.0.0.1"
REDIS_PORT="6379"
REDIS_USERNAME="default"
REDIS_PASSWORD=""
REDIS_CACHE_PREFIX="http-service"
REDIS_CACHE_DURATION=5

# Sentry config
SENTRY_DSN=""
SENTRY_ERROR_TRACING=false
SENTRY_TRACES_SAMPLE_RATE=0.2

# OpenTelemetry config
OTEL_EXPORTER_OTLP_ENDPOINT="localhost:4317"
OTEL_INSECURE_MODE=true

# OAuth 2.0 secrets
OAUTH_PUBLIC_KEY=""
OAUTH_PRIVATE_KEY=""
 ```
