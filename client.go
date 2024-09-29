package kinde

import (
	"context"

	"github.com/axatol/kinde-go/api/apis"
	"github.com/axatol/kinde-go/api/applications"
	"github.com/axatol/kinde-go/api/permissions"
	"github.com/axatol/kinde-go/internal/client"
)

type Client struct {
	client client.Client

	APIs         *apis.Client
	Applications *applications.Client
	Permissions  *permissions.Client
}

func New(ctx context.Context, options *client.ClientOptions) Client {
	client := client.New(ctx, options)

	return Client{
		client:       client,
		APIs:         apis.New(client),
		Applications: applications.New(client),
		Permissions:  permissions.New(client),
	}
}
