# kinde-go

Kinde golang client

## getting started

### prerequisites

1. Create a M2M application and save the domain, client ID, and client secret
2. Authorise the application with the Kinde Management API and allow the relevant scopes

### quickstart

```go
package main

import "github.com/nxt-fwd/kinde-go"

func main() {
  // load config from environment variables
  // - KINDE_DOMAIN
  // - KINDE_AUDIENCE
  // - KINDE_CLIENT_ID
  // - KINDE_CLIENT_SECRET
  client := kinde.New(context.Background(), nil)

  // or override some configuration, the unspecified values are loaded from the
  // environment
  client = kinde.New(
    context.Background(),
    kinde.NewClientOptions().
      WithClientID("foo").
      WithClientSecret("bar").
      WithLogger(someLogger{})
  )
}
```

## development

### testing and linting

The project includes a Makefile with common development tasks:

```bash
# Show all available make targets
make help

# Run integration tests (requires valid API credentials)
make test-e2e

# Generate test coverage report
make coverage

# Run linters
make lint

# Format code
make fmt

# Run tests and linters
make check

# Install development tools (e.g., golangci-lint)
make tools
```

To run integration tests, ensure your `.env` file is properly configured with valid Kinde API credentials.

## todo

- pagination
- rate-limiting
