.PHONY: build run test tidy lint clean

BINARY_NAME=decode-server
BUILD_DIR=.

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v ./...

tidy:
	go mod tidy

lint:
	go vet ./...

clean:
	rm -f $(BINARY_NAME)

deps:
	go mod download