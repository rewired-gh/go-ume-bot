pre:
	go mod tidy
	mkdir -p ./target

test:
	go test ./...

dev:
	go run ./cmd/umed

build: pre
	go build -o ./target ./cmd/umed

build_x64: pre
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o ./target/umed_linux_amd64 ./cmd/umed

.PHONY: pre test dev build build_x64
