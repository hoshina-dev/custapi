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
  └── services/       # Business logic layer
```
