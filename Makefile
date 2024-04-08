run: build/wasm
	go run host.go

build/wasm:
	GOOS=wasip1 GOARCH=wasm go build -o plugin.wasm plugin.go
