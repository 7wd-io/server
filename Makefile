#!make

#!make
ENV_FILE = .env.$(ENV)
# to share env here
include $(ENV_FILE)
# to share env in program
export $(shell sed 's/=.*//' $(ENV_FILE))

DC_CLI = docker compose --env-file=$(ENV_FILE) --profile $(ENV)
PG_DSN = postgres://$(SWD_PG_USER):$(SWD_PG_PASSWORD)@localhost:$(SWD_PG_PORT)/$(SWD_PG_DBNAME)?sslmode=disable

LOCAL_BIN=$(PWD)/scripts/bin
LOCAL_TMP=$(PWD)/scripts/tmp




#.PHONY : run rin-api run-clock run-online

#i:
#	$(info Installing binary dependencies...)
#
#	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

#docker compose -f compose.yml --env-file=.env.dev

default:
	$(info Commands: up,down,migrate)

up:
	$(DC_CLI) up -d

down:
	$(DC_CLI) down

m:
	$(LOCAL_BIN)/migrate create -ext sql -dir migrations -seq  $(NAME)

m-up:
	$(LOCAL_BIN)/migrate -database "postgres://$(SWD_PG_USER):$(SWD_PG_PASSWORD)@localhost:$(SWD_PG_PORT)/$(SWD_PG_DBNAME)?sslmode=disable" -path "migrations" up

m-down:
	$(LOCAL_BIN)/migrate -database "postgres://$(SWD_PG_USER):$(SWD_PG_PASSWORD)@localhost:$(SWD_PG_PORT)/$(SWD_PG_DBNAME)?sslmode=disable" -path "migrations" down

ttr:
	go test ./... -cover -race -vet=all

build:
	go build -C cmd/api -o swd-api
	go build -C cmd/online -o swd-online
	go build -C cmd/clock -o swd-clock

