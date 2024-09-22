//go:build e2e
// +build e2e

package kinde_test

import (
	"context"
	"testing"

	"github.com/axatol/kinde-go"
	"github.com/axatol/kinde-go/internal/e2e"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/require"
)

var e2eClient *kinde.Client

func defaultE2EClientOptions(t *testing.T) *kinde.ClientOptions {
	t.Helper()

	return kinde.NewClientOptions().
		WithDomain(e2e.Domain).
		WithAudience(e2e.Audience).
		WithClientID(e2e.ClientID).
		WithClientSecret(e2e.ClientSecret).
		WithLogger(testutil.NewTestLogger(t))
}

func defaultE2EClient(t *testing.T) *kinde.Client {
	t.Helper()

	if e2eClient != nil {
		return e2eClient
	}

	e2eClient := kinde.New(context.TODO(), defaultE2EClientOptions(t))
	require.NotNil(t, e2eClient)
	return e2eClient
}
