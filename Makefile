include .envrc
MIGRATIONS_PATH = ./cmd/migrate/migrations

test:
	@go test -v ./...

server:
	@go run cmd/server/*.go

client:
	@go run cmd/client/*.go

build:
	@go build -o bin/blaekward-server ./cmd/server/*.go
	@go build -o bin/blaekward-client ./cmd/client/*.go

install:
	@cp -r bin/* /usr/local/bin/

build-windows:
	@GOOS=windows GOARCH=amd64 go build -o bin/blaekward-server.exe ./cmd/server/*.go
	@GOOS=windows GOARCH=amd64 go build -o bin/blaekward-client.exe ./cmd/client/*.go

install-windows:
	@sudo cp -r bin/* /mnt/c/dev 

migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(SERVER_ADDR) up

migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(SERVER_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

docker:
	@docker-compose -f docker-compose.dev.yaml up --build --remove-orphans

docker-prod:
	@docker-compose -f docker-compose.prod.yaml up --build --remove-orphans
