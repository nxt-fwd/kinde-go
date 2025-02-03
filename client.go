package kinde

import (
	"context"

	"github.com/axatol/kinde-go/api/apis"
	"github.com/axatol/kinde-go/api/applications"
	"github.com/axatol/kinde-go/api/organizations"
	"github.com/axatol/kinde-go/api/permissions"
	"github.com/axatol/kinde-go/api/users"
	"github.com/axatol/kinde-go/internal/client"
)

type Client struct {
	client client.Client

	APIs          *apis.Client
	Applications  *applications.Client
	Organizations *organizations.Client
	Permissions   *permissions.Client
	Users         *users.Client
}

func New(ctx context.Context, options *ClientOptions) Client {
	client := client.New(ctx, options.ClientOptions)

	return Client{
		client:        client,
		APIs:         apis.New(client),
		Applications: applications.New(client),
		Organizations: organizations.New(client),
		Permissions:  permissions.New(client),
		Users:        users.New(client),
	}
}
