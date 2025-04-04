test:
	@go test -v ./...

server:
	@go run cmd/server/*.go

client:
	@go run cmd/client/*.go

build-server:
	@go build -o bin/blaekward-server ./cmd/server/*.go

build-client:
	@go build -o bin/blaekward-client ./cmd/client/*.go

install:
	@cp bin/ink-server /usr/local/bin/
	@cp bin/ink-client /usr/local/bin/

build-windows:
	@GOOS=windows GOARCH=amd64 go build -o bin/blaekward-server.exe ./cmd/server/*.go
	@GOOS=windows GOARCH=amd64 go build -o bin/blaekward-client.exe ./cmd/client/*.go

install-windows:
	@sudo cp -r bin/* /mnt/c/dev 
