.PHONY: build run clean install test

BINARY_NAME=sofinco-bot
MAIN_PATH=cmd/bot/main.go

build:
	@echo "Building..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

run:
	@echo "Running..."
	@go run $(MAIN_PATH)

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf data/*.db
	@echo "Clean complete"

install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

test:
	@echo "Running tests..."
	@go test -v ./...

dev:
	@echo "Running in development mode..."
	@go run $(MAIN_PATH)

deploy:
	@echo "Building for production..."
	@CGO_ENABLED=1 go build -ldflags="-s -w" -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Production build complete"
