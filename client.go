package kinde

import (
	"context"

	"github.com/nxt-fwd/kinde-go/api/apis"
	"github.com/nxt-fwd/kinde-go/api/applications"
	"github.com/nxt-fwd/kinde-go/api/identities"
	"github.com/nxt-fwd/kinde-go/api/organizations"
	"github.com/nxt-fwd/kinde-go/api/permissions"
	"github.com/nxt-fwd/kinde-go/api/roles"
	"github.com/nxt-fwd/kinde-go/api/users"
	"github.com/nxt-fwd/kinde-go/api/connections"
	"github.com/nxt-fwd/kinde-go/internal/client"
)

type Client struct {
	client client.Client

	APIs          *apis.Client
	Applications  *applications.Client
	Identities    *identities.Client
	Organizations *organizations.Client
	Permissions   *permissions.Client
	Roles         *roles.Client
	Users         *users.Client
	Connections   *connections.Client
}

func New(ctx context.Context, options *ClientOptions) Client {
	client := client.New(ctx, options.ClientOptions)

	return Client{
		client:        client,
		APIs:          apis.New(client),
		Applications:  applications.New(client),
		Identities:    identities.New(client),
		Organizations: organizations.New(client),
		Permissions:   permissions.New(client),
		Roles:         roles.New(client),
		Users:         users.New(client),
		Connections:   connections.New(client),
	}
}
