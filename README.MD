# Go Sandbox

Learning Go by building a paginated REST API with a React dashboard.
Exploring Go patterns and best practices.

## Tech Stack

- **Go** — Chi router, pgx/v5 driver
- **PostgreSQL** — with migrations and seed data
- **React + TypeScript** — data dashboard
- **Docker** — PostgreSQL container

## What I'm Practicing

- Go project structure (`cmd/` and `internal/` pattern)
- pgx/v5 connection pooling
- Pagination with offset-based queries
- Dynamic query building with filters
- Proper error handling and response patterns
- Environment-based configuration
- Docker Compose for local development

## Quick Start
```bash
git clone https://github.com/vovadyach/go-sandbox
cd go-sandbox
cp .env.example .env
docker compose up db -d
cd api
go mod tidy
go run cmd/migrate/main.go
go run cmd/seed/main.go
go run cmd/server/main.go
```

## API Docs
TBD

## API Endpoints
TBD

## Status
🚧 Work in progress