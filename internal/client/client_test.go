package client_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/nxt-fwd/kinde-go/internal/client"
	"github.com/nxt-fwd/kinde-go/internal/testutil"
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

	client := client.New(context.TODO(), nil)
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

	client := client.New(
		context.TODO(),
		client.NewClientOptions().
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

func TestIncompleteCredentials(t *testing.T) {
	// Test cases with different incomplete credential combinations
	testCases := []struct {
		name         string
		domain       string
		audience     string
		clientID     string
		clientSecret string
		expectError  bool
	}{
		{
			name:         "Missing All Credentials",
			expectError:  true,
		},
		{
			name:         "Missing Client Secret",
			domain:       "https://test.kinde.com",
			audience:     "https://test.kinde.com/api",
			clientID:     "test-client-id",
			expectError:  true,
		},
		{
			name:         "Missing Client ID",
			domain:       "https://test.kinde.com",
			audience:     "https://test.kinde.com/api",
			clientSecret: "test-client-secret",
			expectError:  true,
		},
		{
			name:         "Missing Domain",
			audience:     "https://test.kinde.com/api",
			clientID:     "test-client-id",
			clientSecret: "test-client-secret",
			expectError:  true,
		},
		{
			name:         "Missing Audience",
			domain:       "https://test.kinde.com",
			clientID:     "test-client-id",
			clientSecret: "test-client-secret",
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create client options with test credentials
			options := client.NewClientOptions().
				WithDomain(tc.domain).
				WithAudience(tc.audience).
				WithClientID(tc.clientID).
				WithClientSecret(tc.clientSecret).
				WithLogger(testutil.NewTestLogger(t))

			// Create the client
			client := client.New(context.TODO(), options)

			// Try to make a request
			req, err := client.NewRequest(context.TODO(), "GET", "/test", nil, nil)

			if tc.expectError {
				assert.Error(t, err, "expected error due to incomplete credentials")
				assert.Nil(t, req, "request should be nil when there's an error")
				t.Logf("got expected error: %v\n", err)
			} else {
				assert.NoError(t, err, "unexpected error with complete credentials")
				assert.NotNil(t, req, "request should not be nil when credentials are complete")
			}
		})
	}
}

func TestInvalidCredentials(t *testing.T) {
	// Test cases for invalid credentials
	testCases := []struct {
		name         string
		domain       string
		audience     string
		clientID     string
		clientSecret string
		expectedErr  string
	}{
		{
			name:         "Invalid Client ID and Secret",
			domain:       "https://test.kinde.com",
			audience:     "https://test.kinde.com/api",
			clientID:     "invalid-client-id",
			clientSecret: "invalid-client-secret",
			expectedErr:  "Client authentication failed",
		},
		{
			name:         "Invalid Domain",
			domain:       "https://invalid.kinde.com",
			audience:     "https://test.kinde.com/api",
			clientID:     "test-client-id",
			clientSecret: "test-client-secret",
			expectedErr:  "failed to execute request",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create client options with test credentials
			options := client.NewClientOptions().
				WithDomain(tc.domain).
				WithAudience(tc.audience).
				WithClientID(tc.clientID).
				WithClientSecret(tc.clientSecret).
				WithLogger(testutil.NewTestLogger(t))

			// Create the client
			client := client.New(context.TODO(), options)

			// Try to make a request
			req, err := client.NewRequest(context.TODO(), "GET", "/test", nil, nil)
			if err != nil {
				assert.Contains(t, err.Error(), "invalid client configuration", "expected configuration validation error")
				return
			}

			// If request was created, try to execute it
			err = client.DoRequest(req, nil)
			assert.Error(t, err, "expected error with invalid credentials")
			assert.Contains(t, err.Error(), tc.expectedErr, "error message should indicate invalid credentials")
		})
	}
}
