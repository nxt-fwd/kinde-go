package kinde

import (
	"github.com/nxt-fwd/kinde-go/internal/client"
	"github.com/nxt-fwd/kinde-go/internal/logger"
)

// ClientOptions represents the configuration options for the Kinde client.
type ClientOptions struct {
	*client.ClientOptions
}

// NewClientOptions creates a new ClientOptions instance.
// It loads default values from environment variables:
// - KINDE_DOMAIN
// - KINDE_AUDIENCE
// - KINDE_CLIENT_ID
// - KINDE_CLIENT_SECRET
// - KINDE_SCOPES (space-separated list)
func NewClientOptions() *ClientOptions {
	return &ClientOptions{
		ClientOptions: client.NewClientOptions(),
	}
}

// WithDomain sets the domain for the client options.
func (o *ClientOptions) WithDomain(domain string) *ClientOptions {
	o.ClientOptions.WithDomain(domain)
	return o
}

// WithAudience sets the audience for the client options.
func (o *ClientOptions) WithAudience(audience string) *ClientOptions {
	o.ClientOptions.WithAudience(audience)
	return o
}

// WithClientID sets the client ID for the client options.
func (o *ClientOptions) WithClientID(clientID string) *ClientOptions {
	o.ClientOptions.WithClientID(clientID)
	return o
}

// WithClientSecret sets the client secret for the client options.
func (o *ClientOptions) WithClientSecret(clientSecret string) *ClientOptions {
	o.ClientOptions.WithClientSecret(clientSecret)
	return o
}

// WithScopes sets the scopes for the client options.
func (o *ClientOptions) WithScopes(scopes []string) *ClientOptions {
	o.ClientOptions.WithScopes(scopes)
	return o
}

// WithLogger sets the logger for the client options.
func (o *ClientOptions) WithLogger(logger logger.Logger) *ClientOptions {
	o.ClientOptions.WithLogger(logger)
	return o
}
