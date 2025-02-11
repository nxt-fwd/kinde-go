package client

import (
	"fmt"
	"os"
	"strings"

	"github.com/nxt-fwd/kinde-go/internal/logger"
)

type ClientOptions struct {
	Domain       string
	Audience     string
	ClientID     string
	ClientSecret string
	Scopes       []string
	Logger       logger.Logger
	accessToken  string
}

func NewClientOptions() *ClientOptions {
	return &ClientOptions{
		Domain:       os.Getenv("KINDE_DOMAIN"),
		Audience:     os.Getenv("KINDE_AUDIENCE"),
		ClientID:     os.Getenv("KINDE_CLIENT_ID"),
		ClientSecret: os.Getenv("KINDE_CLIENT_SECRET"),
		Scopes:       strings.Split(os.Getenv("KINDE_SCOPES"), " "),
		Logger:       logger.NoopLogger{},
	}
}

func (o *ClientOptions) WithDomain(domain string) *ClientOptions {
	o.Domain = domain
	return o
}

func (o *ClientOptions) WithAudience(audience string) *ClientOptions {
	o.Audience = audience
	return o
}

func (o *ClientOptions) WithClientID(clientID string) *ClientOptions {
	o.ClientID = clientID
	return o
}

func (o *ClientOptions) WithClientSecret(clientSecret string) *ClientOptions {
	o.ClientSecret = clientSecret
	return o
}

func (o *ClientOptions) WithScopes(scopes []string) *ClientOptions {
	o.Scopes = scopes
	return o
}

func (o *ClientOptions) WithLogger(logger logger.Logger) *ClientOptions {
	o.Logger = logger
	return o
}

func (o *ClientOptions) GetAccessToken() string {
	return o.accessToken
}

func (o *ClientOptions) SetAccessToken(token string) {
	o.accessToken = token
}

// Validate checks if all required options are set
func (o *ClientOptions) Validate() error {
	var missing []string

	if o.Domain == "" {
		missing = append(missing, "domain")
	}
	if o.Audience == "" {
		missing = append(missing, "audience")
	}
	if o.ClientID == "" {
		missing = append(missing, "client_id")
	}
	if o.ClientSecret == "" {
		missing = append(missing, "client_secret")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required Kinde client options: %s", strings.Join(missing, ", "))
	}

	return nil
}
