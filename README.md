# Customer API

A clean, layered REST API for managing users and organizations built with Go, Fiber, and PostgreSQL.

## Architecture

```
sql/
 └── xxx_description.sql

cmd/                    # Application entry point
  └── main.go          # Main application bootstrap

docs/                   # Swagger documentation
  ├── docs.go
  ├── swagger.json
  └── swagger.yaml

internal/              # Private application code
  ├── config/         # Configuration management
  ├── database/             # Database connection & migrations
  ├── handlers/       # HTTP handlers (transport layer)
  ├── middleware/     # Cross-cutting concerns (logging, etc.)
  ├── models/         # Domain models & DTOs
  ├── repositories/   # Data persistence layer
  ├── routes/         # Route definitions
  ├── services/       # Business logic layer
  └── telemetry/      # OpenTelemetry tracing & metrics setup

deploy/otel/            # ADOT collector config + example ECS task definition
```

## Observability

Traces and metrics are instrumented with OpenTelemetry and, in production,
shipped through an ADOT collector sidecar to AWS X-Ray (traces) and
CloudWatch (metrics). See [docs/observability.md](docs/observability.md) for
how it's wired up, configuration env vars, and how to deploy the ADOT
sidecar on ECS.
