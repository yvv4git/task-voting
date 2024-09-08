# Environment file
ENV_FILE := .env
ifneq (,$(wildcard ${ENV_FILE}))
    include ${ENV_FILE}
endif

# Default values
GO ?= go

# Go binary directory
GO_BIN_DIR := $(shell ${GO} env GOBIN)
ifeq (,${GO_BIN_DIR})
    GO_BIN_DIR := ${GOPATH}/bin
endif

# Migration tool and paths
MIGRATE_BIN := ${GO_BIN_DIR}/goose
MIGRATIONS_VOTING_APP := ./migrations/voting

# Database connection strings
VOTING_PG_DSN ?= "user=${VOTING_PG_USER} password=${VOTING_PG_PASSWORD} dbname=${VOTING_PG_DB} host=${VOTING_PG_HOST} port=${VOTING_PG_PORT} sslmode=disable"

# Test environment
DOCKER_COMPOSE_DEV_ENV ?= docker-compose.yaml
DOCKER_COMPOSE_DEV_ENV_SERVICES ?= voting-db
DOCKER_COMPOSE_RUN?=docker compose -f ${DOCKER_COMPOSE_DEV_ENV} --project-name=test-task-voting --env-file=${ENV_FILE}


# Dev environment
.PHONY: dev-env-up
dev-env-up:
	${DOCKER_COMPOSE_RUN} up -d ${DOCKER_COMPOSE_TEST_ENV_SERVICES}

.PHONY: dev-env-down
dev-env-down:
	${DOCKER_COMPOSE_RUN} down

.PHONY: dev-env-clear
dev-env-clear: dev-env-down
	# Remove containers and volumes
	${DOCKER_COMPOSE_RUN} down -v

# Migration targets
.PHONY: migrate-voting-up
migrate-voting-up:
	${MIGRATE_BIN} -dir=$(MIGRATIONS_VOTING_APP) postgres $(VOTING_PG_DSN) up

.PHONY: migrate-voting-reset
migrate-voting-reset:
	${MIGRATE_BIN} -dir=$(MIGRATIONS_VOTING_APP) postgres $(VOTING_PG_DSN) reset

run_application:
	go run main.go -c config.toml user