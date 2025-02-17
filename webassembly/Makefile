up:
	cd cmd/server && go run main.go
.PHONY: up

build:
	cd cmd/wasm && GOOS=js GOARCH=wasm go build -o  ../../assets/json.wasm
.PHONY: build