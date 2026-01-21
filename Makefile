APP_NAME := gomonitor
CMD_DIR  := ./cmd/api

.PHONY: run test test-cover build clean

# Run app
run:
	ENVIRONMENT=development go run $(CMD_DIR)

# Run tests
test:
	ENVIRONMENT=test go test ./...

# Run tests with coverage output.
test-cover:
	ENVIRONMENT=test go test ./... -coverprofile=coverage.out