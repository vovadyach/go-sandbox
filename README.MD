# Go REST API

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

## API Endpoints
```
GET /health                          — Health check
GET /api/users?page=1&limit=20       — List users (paginated)
GET /api/users?role=admin            — Filter by role
GET /api/users?status=active         — Filter by status
GET /api/users/:id                   — User detail
GET /api/users/:id/posts             — User's posts (paginated)
GET /api/posts?page=1&limit=20       — List posts (paginated)
GET /api/posts?status=published      — Filter by status
```

## Status
🚧 Work in progress
