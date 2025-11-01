
TARGET = generate-locations

run:
	go run tools/${TARGET}.go

module:
	rm -f go.mod go.sum
	go mod init ${TARGET}
	go mod tidy


clean:
	rm -rf bin
	rm -f go.mod go.sum
