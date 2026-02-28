# Go Sandbox

Building paginated REST API with dynamic filtering and sorting for user and post management

## Tech Stack

- **Go** — Chi router, pgx/v5 driver
- **PostgreSQL** — with migrations and seed data
- **React + TypeScript** — data dashboard
- **Docker** — PostgreSQL container
- **Goose** — database migrations
- **Make** — project commands

## Quick Start
```bash
git clone https://github.com/vovadyach/go-sandbox
cd go-sandbox
cp .env.example .env
docker compose up db -d
make migrate-up
cd api
go mod tidy
go run cmd/seed/main.go
go run cmd/server/main.go
```

## Available Commands
```bash
make migrate-up                        # Run migrations
make migrate-down                      # Rollback migrations
make migrate-status                    # Check migration status
make migrate-create name=add_table     # Create new migration
```

## API Docs
TBD

## Status
🚧 Work in progress
