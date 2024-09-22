package clientcredentials_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/axatol/kinde-go/internal/clientcredentials"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuth2TransportMissingToken(t *testing.T) {
	testConfig := testutil.DefaultTestServerConfig()
	testServer := testutil.NewTestServer(t, testConfig)

	transport := clientcredentials.OAuth2Transport{
		Domain:       testServer.Server.URL,
		Audience:     testConfig.Audience,
		ClientID:     testConfig.ClientID,
		ClientSecret: testConfig.ClientSecret,
		Transport:    http.DefaultTransport,
		Logger:       testutil.NewTestLogger(t),
		Token:        "",
		Expiry:       time.Now().Add(time.Minute),
	}

	req, err := http.NewRequest(http.MethodGet, testServer.Server.URL+"/hello", nil)
	assert.NoError(t, err)
	assert.NotNil(t, req)

	res, err := transport.RoundTrip(req)
	assert.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPost, "/oauth2/token"))
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/hello"))
	assert.Equal(t, testConfig.AccessToken, transport.Token)
}

func TestOAuth2TransportValidToken(t *testing.T) {
	testConfig := testutil.DefaultTestServerConfig()
	testServer := testutil.NewTestServer(t, testConfig)

	transport := clientcredentials.OAuth2Transport{
		Domain:       testServer.Server.URL,
		Audience:     testConfig.Audience,
		ClientID:     testConfig.ClientID,
		ClientSecret: testConfig.ClientSecret,
		Transport:    http.DefaultTransport,
		Logger:       testutil.NewTestLogger(t),
		Token:        testConfig.AccessToken,
		Expiry:       time.Now().Add(time.Minute),
	}

	req, err := http.NewRequest(http.MethodGet, testServer.Server.URL+"/hello", nil)
	assert.NoError(t, err)
	assert.NotNil(t, req)

	res, err := transport.RoundTrip(req)
	assert.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, 0, testServer.CallCount.Get(http.MethodPost, "/oauth2/token"))
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/hello"))
}

func TestOauth2TransportExpiredToken(t *testing.T) {
	testConfig := testutil.DefaultTestServerConfig()
	testServer := testutil.NewTestServer(t, testConfig)

	transport := clientcredentials.OAuth2Transport{
		Domain:       testServer.Server.URL,
		Audience:     testConfig.Audience,
		ClientID:     testConfig.ClientID,
		ClientSecret: testConfig.ClientSecret,
		Transport:    http.DefaultTransport,
		Logger:       testutil.NewTestLogger(t),
		Token:        "expired token",
		Expiry:       time.Now().Add(-time.Minute),
	}

	req, err := http.NewRequest(http.MethodGet, testServer.Server.URL+"/hello", nil)
	assert.NoError(t, err)
	assert.NotNil(t, req)

	res, err := transport.RoundTrip(req)
	assert.NoError(t, err)
	require.NotNil(t, res)

	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPost, "/oauth2/token"))
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/hello"))
}
