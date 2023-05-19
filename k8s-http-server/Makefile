
TARGET = k8s-http-server

build:
	go build -o bin/${TARGET} cmd/${TARGET}.go

.PHONY: test
test:
	HOSTNAME="lobo-codes-abcdef1234-vwxyz" NODE_NAME="work-macbook" go run cmd/${TARGET}.go -port=":3000" \
		-ssl-port=":3001" -cassandra-db="10.1.1.42" -os-release="test/os-release" -go-version="test/golang_version.txt"

module:
	rm -f go.mod go.sum
	go mod init ${TARGET}
	go mod tidy

.PHONY: docker
docker:
	docker build -t localhost:32000/lobo-codes:$(tag) -f docker/Dockerfile .
	docker tag localhost:32000/lobo-codes:$(tag) localhost:32000/lobo-codes:latest
	docker push localhost:32000/lobo-codes:$(tag)
	docker push localhost:32000/lobo-codes:latest

rollout:
	microk8s kubectl rollout restart  -f ./k8s/lobo-codes-deployment.yaml

transfer:
	go run cmd/transfer-sqlite-cassandra.go -sqliteDb=requests-db/$(name).db -cassandraDb=10.1.1.42 -cassandraTbl=$(name)

k8s-test:
	docker build -t localhost:32000/k8s-test:latest -f docker/Dockerfile.k8s-test .
	docker push localhost:32000/k8s-test:latest
	microk8s kubectl run --rm -i k8s-test --image=localhost:32000/k8s-test

clean:
	rm -rf bin
	rm -f go.mod go.sum
