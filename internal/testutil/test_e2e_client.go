//go:build e2e
// +build e2e

package testutil

import (
	"context"
	"testing"

	"github.com/nxt-fwd/kinde-go/internal/client"
	"github.com/nxt-fwd/kinde-go/internal/e2e"
	"github.com/stretchr/testify/require"
)

var e2eClient client.Client

func DefaultE2EClientOptions(t *testing.T) *client.ClientOptions {
	t.Helper()

	return client.NewClientOptions().
		WithDomain(e2e.Domain).
		WithAudience(e2e.Audience).
		WithClientID(e2e.ClientID).
		WithClientSecret(e2e.ClientSecret).
		WithLogger(NewTestLogger(t))
}

func DefaultE2EClient(t *testing.T) client.Client {
	t.Helper()

	if e2eClient != nil {
		return e2eClient
	}

	e2eClient := client.New(context.TODO(), DefaultE2EClientOptions(t))
	require.NotNil(t, e2eClient)
	return e2eClient
}
