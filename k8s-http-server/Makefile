
TARGET = k8s-http-server

build:
	go build -o bin/${TARGET} cmd/${TARGET}.go

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
