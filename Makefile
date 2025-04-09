test:
	@go test -v ./...

server:
	@go run cmd/controller/*.go

client:
	@go run cmd/gateway/*.go

build:
	@go build -o bin/blaekward-controller ./cmd/controller/*.go
	@go build -o bin/blaekward-gateway ./cmd/gateway/*.go

install:
	@cp bin/ink-server /usr/local/bin/
	@cp bin/ink-client /usr/local/bin/

build-windows:
	@GOOS=windows GOARCH=amd64 go build -o bin/blaekward-controller.exe ./cmd/controller/*.go
	@GOOS=windows GOARCH=amd64 go build -o bin/blaekward-gateway.exe ./cmd/gateway/*.go

install-windows:
	@sudo cp -r bin/* /mnt/c/dev 
