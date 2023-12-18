#!make

ENV_FILE=.env.$(ENV)
DC_FILE=docker-compose.$(ENV).yml
LOCAL_BIN=$(PWD)/scripts/bin
LOCAL_TMP=$(PWD)/scripts/tmp

# to share env here
include $(ENV_FILE)
# to share env in program
export $(shell sed 's/=.*//' $(ENV_FILE))

i:
	$(info Installing binary dependencies...)

	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

dc-up:
	docker compose --file=$(DC_FILE) --env-file=$(ENV_FILE) up -d

dc-down:
	docker compose --file=$(DC_FILE) --env-file=$(ENV_FILE) down

m:
	$(LOCAL_BIN)/migrate create -ext sql -dir migrations -seq  $(NAME)

m-up:
	$(LOCAL_BIN)/migrate -database "postgres://$(ENV):$(ENV)@localhost:$(SWD_PG_PORT)/$(SWD_PG_DBNAME)?sslmode=disable" -path "migrations" up

m-down:
	$(LOCAL_BIN)/migrate -database "postgres://$(ENV):$(ENV)@localhost:$(SWD_PG_PORT)/$(SWD_PG_DBNAME)?sslmode=disable" -path "migrations" down

ttr:
	go test ./... -cover -race -vet=all
