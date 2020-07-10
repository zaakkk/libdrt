# This is comment

task:
	GOOS=js
	GOARCH=wasm
	sudo go build -o main.wasm sample.go
