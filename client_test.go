package kinde_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/axatol/kinde-go"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientConfigWithEnv(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	t.Cleanup(testServer.Server.Close)

	t.Setenv("KINDE_DOMAIN", testServer.Server.URL)
	t.Setenv("KINDE_AUDIENCE", testServer.Config.Audience)
	t.Setenv("KINDE_CLIENT_ID", testServer.Config.ClientID)
	t.Setenv("KINDE_CLIENT_SECRET", testServer.Config.ClientSecret)

	client := kinde.New(context.TODO(), nil)
	require.NotNil(t, client)

	req, err := client.NewRequest(context.TODO(), http.MethodGet, "/hello", nil, nil)
	assert.NoError(t, err)
	require.NotNil(t, req)

	var result map[string]string
	err = client.DoRequest(req, &result)
	assert.NoError(t, err)
	require.NotNil(t, result)

	code, ok := result["code"]
	assert.True(t, ok)
	assert.Equal(t, code, "OK")
}

func TestClientConfigWithOptions(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	t.Cleanup(testServer.Server.Close)

	client := kinde.New(
		context.TODO(),
		kinde.NewClientOptions().
			WithDomain(testServer.Server.URL).
			WithAudience(testServer.Config.Audience).
			WithClientID(testServer.Config.ClientID).
			WithClientSecret(testServer.Config.ClientSecret),
	)
	require.NotNil(t, client)

	req, err := client.NewRequest(context.TODO(), http.MethodGet, "/hello", nil, nil)
	assert.NoError(t, err)
	require.NotNil(t, req)

	var result map[string]string
	err = client.DoRequest(req, &result)
	assert.NoError(t, err)
	require.NotNil(t, result)

	code, ok := result["code"]
	assert.True(t, ok)
	assert.Equal(t, code, "OK")
}
