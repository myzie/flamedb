
DEFAULT_DB_USER ?= curtis

DB_NAME ?= flame
DB_USER ?= flame
DB_PASS ?= flame
DB_HOST ?= 127.0.0.1

.PHONY: all
all:
	go get -v ./...
	go build -v .

.PHONY: test
test:
	go test -v -cover -coverprofile=coverage.out ./...

.PHONY: coverage
coverage:
	go tool cover -html=coverage.out

.PHONY: generate
generate:
	swagger generate server -A flamedb -P models.Principal -f ./swagger.yaml

.PHONY: mocks
mocks:
	mkdir -p ./database/mock_database
	go generate ./database

.PHONY: install_swagger
install_swagger:
	brew tap go-swagger/go-swagger
	brew install go-swagger

.PHONY: db
db:
	psql -U $(DEFAULT_DB_USER) -d postgres \
		-c "CREATE ROLE $(DB_USER) WITH LOGIN PASSWORD '$(DB_PASS)';"
	psql -U $(DEFAULT_DB_USER) -d postgres \
		-c "CREATE DATABASE $(DB_NAME) OWNER $(DB_USER);"

.PHONY: run
run:
	DB_NAME=$(DB_NAME) DB_USER=$(DB_USER)     \
	DB_PASSWORD=$(DB_PASS) DB_HOST=$(DB_HOST) \
	./flamedb
