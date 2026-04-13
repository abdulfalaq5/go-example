.PHONY: run build tidy lint test

APP_NAME=go-example
BUILD_DIR=./bin
MAIN=./cmd/api/main.go

## run: start the development server
run:
	go run $(MAIN)

## build: compile to binary
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)
	@echo "Binary: $(BUILD_DIR)/$(APP_NAME)"

## tidy: clean and sync dependencies
tidy:
	go mod tidy

## test: run all tests
test:
	go test ./... -v -race -count=1

## lint: run golangci-lint (must be installed)
lint:
	golangci-lint run ./...

## help: print available targets
help:
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
