
TARGET = k8s-http-server

DOCKER = $(shell which docker)

build:
	go build -o bin/${TARGET} cmd/${TARGET}.go

debug:
	go build -gcflags "all=-N -l" -o bin/${TARGET} cmd/${TARGET}.go

.PHONY: test
test:
	HOSTNAME="lobo-codes-abcdef1234-vwxyz" NODE_NAME="work-macbook" go run cmd/${TARGET}.go \
		-os-release="test/os-release" -go-version="test/golang_version.txt"

test-debug:
	HOSTNAME="lobo-codes-abcdef1234-vwxyz" NODE_NAME="work-macbook" \
		dlv --listen=:2345 --headless=true --api-version=2 --accept-multiclient exec ./bin/${TARGET} \
		-- -os-release="test/os-release" -go-version="test/golang_version.txt"

module:
	rm -f go.mod go.sum
	go mod init ${TARGET}
	go mod tidy

.PHONY: docker
docker:
	${DOCKER} build -t localhost:32000/lobo-codes:$(tag) -f docker/Dockerfile .
	${DOCKER} tag localhost:32000/lobo-codes:$(tag) localhost:32000/lobo-codes:latest
	${DOCKER} push localhost:32000/lobo-codes:$(tag)
	${DOCKER} push localhost:32000/lobo-codes:latest

test-docker:
	${DOCKER} run --rm -p 3000:3000 -e HOSTNAME="lobo-codes-docker1234-abcde" -e NODE_NAME="docker-node" localhost:32000/lobo-codes:latest

local-docker:
	${DOCKER} build -t lobo-codes:$(tag) -f docker/Dockerfile .
	${DOCKER} tag lobo-codes:$(tag) lobo-codes:latest

test-local-docker:
	${DOCKER} run --rm -p 3000:3000 -e HOSTNAME="lobo-codes-docker1234-abcde" -e NODE_NAME="docker-node" lobo-codes:latest

deploy-local-docker:
	${DOCKER} run --rm -p 3000:80 -e HOSTNAME="lobo-codes-docker1234-abcde" -e NODE_NAME="docker-node" lobo-codes:latest

apply:
	microk8s kubectl apply -f ./k8s/lobo-codes-deployment.yaml

delete:
	microk8s kubectl delete -f ./k8s/lobo-codes-deployment.yaml

rollout:
	microk8s kubectl rollout restart -f ./k8s/lobo-codes-deployment.yaml

transfer:
	go run cmd/transfer-sqlite-cassandra.go -sqliteDb=requests-db/$(name).db -cassandraDb=10.1.1.42 -cassandraTbl=$(name)

build-transfer2:
	go build -o bin/transfer-cassandra-rqlite cmd/transfer-cassandra-rqlite.go

run-transfer2:
	./bin/transfer-cassandra-rqlite -cassandraTbl=$(name)

k8s-test:
	docker build -t localhost:32000/k8s-test:latest -f docker/Dockerfile.k8s-test .
	docker push localhost:32000/k8s-test:latest
	microk8s kubectl run --rm -i k8s-test --image=localhost:32000/k8s-test

clean:
	rm -rf bin
	rm -f go.mod go.sum
