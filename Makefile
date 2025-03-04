#!make

#!make
ENV_FILE = .env.$(ENV)
# to share env here
include $(ENV_FILE)
# to share env in program
export $(shell sed 's/=.*//' $(ENV_FILE))

DC_CLI = docker compose --env-file=$(ENV_FILE) --profile $(ENV)
MIGRATE_CLI = $(LOCAL_BIN)/migrate
PG_DSN = postgres://$(SWD_PG_USER):$(SWD_PG_PASSWORD)@localhost:$(SWD_PG_PORT)/$(SWD_PG_DBNAME)?sslmode=disable

LOCAL_BIN=$(PWD)/scripts/bin
LOCAL_TMP=$(PWD)/scripts/tmp
DEPLOY_BIN=/usr/local/bin


#.PHONY : run rin-api run-clock run-online

i:
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

#docker compose -f compose.yml --env-file=.env.dev

default:
	$(info Commands: up,down,migrate)

up:
	$(DC_CLI) up -d

down:
	$(DC_CLI) down

m:
	$(MIGRATE_CLI) create -ext sql -dir migrations -seq  $(NAME)

m-up:
	$(MIGRATE_CLI) -database "postgres://$(SWD_PG_USER):$(SWD_PG_PASSWORD)@localhost:$(SWD_PG_PORT)/$(SWD_PG_DBNAME)?sslmode=disable" -path "migrations" up

m-down:
	$(MIGRATE_CLI) -database "postgres://$(SWD_PG_USER):$(SWD_PG_PASSWORD)@localhost:$(SWD_PG_PORT)/$(SWD_PG_DBNAME)?sslmode=disable" -path "migrations" down

ttr:
	go test ./... -cover -race -vet=all

build:
	go build -C cmd/api -o $DEPLOY_BIN/swd-api
	go build -C cmd/online -o $DEPLOY_BIN/swd-online
	go build -C cmd/clock -o $DEPLOY_BIN/swd-clock

# https://7wd.io.local/welcome