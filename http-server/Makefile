
TARGET = http-server

run-prod:
	./bin/${TARGET} -log=logs.txt

run-stage:
	./bin/${TARGET} -port=:3000 -ssl-port=:3001

run-test:
	./bin/${TARGET} -port=:3000 -ssl-port= -html-dir=static \
		-geodb=../ip2location-db/IP2LOCATION-LITE-DB11.IPV6.BIN -requestsdb=requests.db

build:
	go build -o bin/${TARGET} cmd/${TARGET}.go

module:
	rm -f go.mod go.sum
	go mod init ${TARGET}
	go mod tidy

clean:
	rm -rf bin
	rm -f go.mod go.sum
