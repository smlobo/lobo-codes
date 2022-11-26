
TARGET = http-server

run:
	./bin/${TARGET} :3000 static ../ip2location-db/IP2LOCATION-LITE-DB11.IPV6.BIN requests.db

build:
	go build -o bin/${TARGET} cmd/${TARGET}.go

module:
	rm -f go.mod go.sum
	go mod init ${TARGET}
	go mod tidy

clean:
	rm -rf bin
	rm -f go.mod go.sum
