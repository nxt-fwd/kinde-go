package kinde

import (
	"os"
	"strings"

	"github.com/axatol/kinde-go/internal/logger"
)

type ClientOptions struct {
	Domain       string
	Audience     string
	ClientID     string
	ClientSecret string
	Scopes       []string
	Logger       logger.Logger
}

func NewClientOptions() *ClientOptions {
	return &ClientOptions{
		Domain:       os.Getenv("KINDE_DOMAIN"),
		Audience:     os.Getenv("KINDE_AUDIENCE"),
		ClientID:     os.Getenv("KINDE_CLIENT_ID"),
		ClientSecret: os.Getenv("KINDE_CLIENT_SECRET"),
		Scopes:       strings.Fields(os.Getenv("KINDE_SCOPES")),
		Logger:       &logger.NoopLogger{},
	}
}

func (c *ClientOptions) WithDomain(value string) *ClientOptions {
	c.Domain = value
	return c
}

func (c *ClientOptions) WithAudience(value string) *ClientOptions {
	c.Audience = value
	return c
}

func (c *ClientOptions) WithClientID(value string) *ClientOptions {
	c.ClientID = value
	return c
}

func (c *ClientOptions) WithClientSecret(value string) *ClientOptions {
	c.ClientSecret = value
	return c
}

func (c *ClientOptions) WithScopes(value []string) *ClientOptions {
	c.Scopes = value
	return c
}

func (c *ClientOptions) WithLogger(value logger.Logger) *ClientOptions {
	c.Logger = value
	return c
}
