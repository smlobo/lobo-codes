
TARGET = kd-tree

build:
	GOOS=js GOARCH=wasm go build -o ${TARGET}.wasm main.go

module:
	rm -f go.mod go.sum
	go mod init ${TARGET}
	go mod tidy
	cp "$(shell go env GOROOT)/lib/wasm/wasm_exec.js" .

install: ${TARGET}.html ${TARGET}.wasm wasm_exec.js
	mkdir -p dist/${TARGET}
	cp ${TARGET}.html dist/${TARGET}/index.html
	cp ${TARGET}.wasm dist/${TARGET}/
	cp wasm_exec.js dist/${TARGET}/

clean:
	rm -rf *.js *.wasm
	rm -f go.mod go.sum
	rm -rf dist
