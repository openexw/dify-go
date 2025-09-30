.PHONY: test coverage

# Run all tests

test:
	go test -v ./...

# Generate test coverage report

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out