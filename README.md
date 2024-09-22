# kinde-go

Kinde golang client

## getting started

### prerequisites

1. Create a M2M application and save the domain, client ID, and client secret
2. Authorise the application with the Kinde Management API and allow the relevant scopes

### quickstart

```go
package main

import "github.com/axatol/kinde-go"

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

## todo

- pagination
- rate-limiting
