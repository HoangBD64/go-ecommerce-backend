include .env

GOOSE_DBSTRING ?= $(STR_MYSQL)"root:root1234@tcp(127.0.0.1:33306)/shopdevgo"
GOOSE_MIGRATION_DIR ?= sql/schema
GOOSE_DRIVER ?= mysql

# Tên của ứng dụng của bạn
APP_NAME := server

# test

all:
	@echo "The message is: $(MESSAGE)"
	@echo "The count is: $(COUNT)"

print_vars:
	@echo "MESSAGE (from make): $(MESSAGE)"
	@echo "COUNT (from make): $(COUNT)"

# Chạy ứng dụng

docker_build:
	docker-compose up -d --build
	docker-compose ps
docker_down:
	docker-compose -f environment/docker-compose-dev.yml down

dev:
	go run ./cmd/$(APP_NAME)

docker_up:
	docker-compose -f environment/docker-compose-dev.yml up

up_by_one:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) up-by-one
# create new a migration
create_migration:
	@goose -dir=$(GOOSE_MIGRATION_DIR) create $(name) sql
upse:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) up
downse:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) down
resetse:
	@GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir=$(GOOSE_MIGRATION_DIR) reset

sqlgen:
	sqlc generate

swag:
	swag init -g ./cmd/server/main.go -o ./cmd/swag/docs

.PHONY: dev downse upse resetse docker_build docker_stop docker_up swag

.PHONY: air
