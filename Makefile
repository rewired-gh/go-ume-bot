pre:
	go mod tidy
	mkdir -p ./target

dev:
	go run ./cmd/umed

build: pre
	go build -o ./target ./cmd/umed

build_x64: pre
	GOOS=linux GOARCH=amd64 go build -o ./target/umed_linux_amd64 ./cmd/umed

.PHONY: pre dev build build_x64