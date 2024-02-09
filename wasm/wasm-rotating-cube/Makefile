TARGET = wasm-rotating-cube

all: ${TARGET}

wasm-rotating-cube: output.wasm main.go
	go build -o wasm-rotating-cube main.go

output.wasm: bundle.go wasm_exec.js
	GOOS=js GOARCH=wasm go build -o bundle.wasm bundle.go

run: output.wasm wasm-rotating-cube
	./wasm-rotating-cube

wasm_exec.js:
	cp $(GOROOT)/misc/wasm/wasm_exec.js .

module:
	rm -f go.mod go.sum
	go mod init ${TARGET}
	go mod tidy

clean:
	rm -f wasm-rotating-cube wasm_exec.js *.wasm
	rm -f go.mod go.sum
