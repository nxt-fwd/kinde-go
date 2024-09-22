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

func defaultScenario(t *testing.T) (*kinde.Client, *testutil.TestServer) {
	t.Helper()

	config := testutil.DefaultTestServerConfig()
	server := testutil.NewTestServer(t, config)

	client := kinde.New(
		context.TODO(),
		kinde.NewClientOptions().
			WithDomain(server.Server.URL).
			WithAudience(config.Audience).
			WithClientID(config.ClientID).
			WithClientSecret(config.ClientSecret).
			WithLogger(testutil.NewTestLogger(t)),
	)
	require.NotNil(t, client)

	return client, server
}

func TestClientConfigWithEnv(t *testing.T) {
	testConfig := testutil.DefaultTestServerConfig()
	testServer := testutil.NewTestServer(t, testConfig)
	t.Cleanup(testServer.Server.Close)

	t.Setenv("KINDE_DOMAIN", testServer.Server.URL)
	t.Setenv("KINDE_AUDIENCE", testConfig.Audience)
	t.Setenv("KINDE_CLIENT_ID", testConfig.ClientID)
	t.Setenv("KINDE_CLIENT_SECRET", testConfig.ClientSecret)

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
	testConfig := testutil.DefaultTestServerConfig()
	testServer := testutil.NewTestServer(t, testConfig)
	t.Cleanup(testServer.Server.Close)

	client := kinde.New(
		context.TODO(),
		kinde.NewClientOptions().
			WithDomain(testServer.Server.URL).
			WithAudience(testConfig.Audience).
			WithClientID(testConfig.ClientID).
			WithClientSecret(testConfig.ClientSecret),
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
