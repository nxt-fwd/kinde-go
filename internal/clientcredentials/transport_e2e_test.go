//go:build e2e
// +build e2e

package clientcredentials_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/axatol/kinde-go/internal/clientcredentials"
	"github.com/axatol/kinde-go/internal/e2e"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestE2ETransportRefresh(t *testing.T) {
	transport := clientcredentials.OAuth2Transport{
		Domain:       e2e.Domain,
		Audience:     e2e.Audience,
		ClientID:     e2e.ClientID,
		ClientSecret: e2e.ClientSecret,
		Transport:    http.DefaultTransport,
		Logger:       testutil.NewTestLogger(t),
	}

	token, err := transport.GetToken(context.TODO())
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
