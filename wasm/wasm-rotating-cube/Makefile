TARGET = wasm-rotating-cube

GOROOT := $(shell go env GOROOT)

all: ${TARGET}

${TARGET}: bundle.wasm main.go
	go build -o $@ main.go

bundle.wasm: bundle.go wasm_exec.js
	GOOS=js GOARCH=wasm go build -o $@ bundle.go

run: ${TARGET} index.html
	python3 -m http.server 9000

dist: ${TARGET} index.html
	mkdir -p $@/${TARGET}
	cp index.html $@/${TARGET}
	cp wasm_exec.js $@/${TARGET}
	cp bundle.wasm $@/${TARGET}
	cp $< $@/${TARGET}

wasm_exec.js:
	cp $(GOROOT)/misc/wasm/wasm_exec.js .

module:
	rm -f go.mod go.sum
	go mod init ${TARGET}
	go mod tidy

clean:
	rm -f ${TARGET} wasm_exec.js *.wasm
	rm -f go.mod go.sum
