
TARGET = k8s-http-server

build:
	go build -o bin/${TARGET} cmd/${TARGET}.go

.PHONY: test
test:
	go run cmd/${TARGET}.go -port=":3000" -ssl-port=":3001" -cassandra-db="10.1.1.42" -os-release="test/os-release"

module:
	rm -f go.mod go.sum
	go mod init ${TARGET}
	go mod tidy

docker:
	docker build -t localhost:32000/lobo-codes:$(tag) .
	docker tag localhost:32000/lobo-codes:$(tag) localhost:32000/lobo-codes:latest
	docker push localhost:32000/lobo-codes:$(tag)
	docker push localhost:32000/lobo-codes:latest

rollout:
	microk8s kubectl rollout restart  -f ./k8s/lobo-codes-deployment.yaml

transfer:
	go run cmd/transfer-sqlite-cassandra.go -sqliteDb=requests-db/$(name).db -cassandraDb=10.1.1.42 -cassandraTbl=$(name)

clean:
	rm -rf bin
	rm -f go.mod go.sum
