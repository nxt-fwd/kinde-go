.PHONY: help test test-unit test-e2e lint lint-fix fmt check clean tools coverage release

# Default target
help:
	@echo "Available targets:"
	@echo "  help       - Show this help message"
	@echo "  test-e2e   - Run end-to-end tests only"
	@echo "  coverage   - Generate test coverage report"
	@echo "  lint       - Run linters"
	@echo "  lint-fix   - Run linters with auto-fix enabled"
	@echo "  fmt        - Format code using gofmt"
	@echo "  check      - Run tests and linters"
	@echo "  clean      - Remove build artifacts"
	@echo "  release    - Create a new release (usage: make release VERSION=v1.2.3)"

# Run e2e tests only
test-e2e:
	@echo "Running e2e tests..."
	go test -v -tags e2e ./api/...

# Generate test coverage
coverage:
	@echo "Generating test coverage report..."
	go test -v ./... -short -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run golangci-lint
lint: tools
	@echo "Running linters..."
	golangci-lint run

# Run golangci-lint with auto-fix enabled
lint-fix: tools
	@echo "Running linters with auto-fix..."
	golangci-lint run --fix

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run tests and linting
check: test-e2e lint

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf dist/ coverage.out coverage.html
	go clean -testcache

# Release a new version
release:
	@if [ "$(VERSION)" = "" ]; then \
		echo "Error: VERSION is required. Usage: make release VERSION=v1.2.3"; \
		exit 1; \
	fi
	@if ! echo "$(VERSION)" | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' > /dev/null; then \
		echo "Error: VERSION must be in format v1.2.3"; \
		exit 1; \
	fi
	@if [ "$$(git rev-parse --abbrev-ref HEAD)" != "master" ]; then \
		echo "Error: Releases must be created from the main branch"; \
		exit 1; \
	fi
	@if git rev-parse "$(VERSION)" >/dev/null 2>&1; then \
		echo "Error: Tag $(VERSION) already exists"; \
		exit 1; \
	fi
	@echo "Creating release $(VERSION)..."
	@echo "Running tests..."
	@make test-e2e
	@echo "Running linter..."
	@make lint
	@echo "Creating git tag..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "Pushing tag to remote..."
	@git push origin $(VERSION)
	@echo "Release $(VERSION) created successfully!"
