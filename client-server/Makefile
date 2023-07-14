
tidy:
	#go mod init client-server
	go mod tidy

server:
	go run cmd/server/server.go

client:
	go run cmd/client/client.go

.PHONY: test
test:
	go test ./test

build:
	go build -o bin/server cmd/server/server.go
	go build -o bin/client cmd/client/client.go

clean:
	rm -f go.sum
	rm -f server client
