include .envrc

MIGRATIONS_PATH = ./cmd/migrate/migrations
SERVER_ADDR = postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost/$(POSTGRES_DB)?sslmode=disable

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

migrate:
	@read -p "Enter the sequence name: " SEQ; \
	migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $${SEQ}

migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(SERVER_ADDR) up

migrate-down:
	@read -p "Number of migrations you want to rollback (default: 1): " NUM; NUM=$${NUM:-1}; \
	migrate -path=$(MIGRATIONS_PATH) -database=$(SERVER_ADDR) down $${NUM}

migrate-force: 
	@read -p "Enter the version to force: " VERSION; \
	@migrate -path=$(MIGRATIONS_PATH) -database=$(SERVER_ADDR) force $${VERSION}

docker:
	@docker-compose -f docker-compose.dev.yaml up --build --remove-orphans

docker-prod:
	@docker-compose -f docker-compose.prod.yaml up --build --remove-orphans
