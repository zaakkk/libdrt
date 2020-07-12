# This is comment

build:
	set GOOS=js
	set GOARCH=wasm
	go build -o main.wasm sample.go

test:
	go run server.go
