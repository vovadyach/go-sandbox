include .env

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

MIGRATIONS_DIR ?= api/db/migrations
GOOSE_DRIVER   ?= postgres


migrate-up:
	goose -dir "$(MIGRATIONS_DIR)" $(GOOSE_DRIVER) "$(DB_URL)" up


migrate-down:
	goose -dir "$(MIGRATIONS_DIR)" $(GOOSE_DRIVER) "$(DB_URL)" down

migrate-status:
	goose -dir "$(MIGRATIONS_DIR)" $(GOOSE_DRIVER) "$(DB_URL)" status

migrate-create:
	@test -n "$(name)" || (echo "Usage: make migrate-create name=<name>" && exit 1)
	goose -dir "$(MIGRATIONS_DIR)" create "$(name)" sql
