//go:build e2e
// +build e2e

package kinde_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/axatol/kinde-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EGetAPIs(t *testing.T) {
	client := defaultE2EClient(t)
	res, err := client.GetAPIs(context.TODO())
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestE2EGetAPI(t *testing.T) {
	client := defaultE2EClient(t)
	res, err := client.GetAPIs(context.TODO())
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestE2ECreateGetDeleteAPI(t *testing.T) {
	client := defaultE2EClient(t)
	// client, server := defaultScenario(t)
	// server.HandleAuthenticated(t, http.MethodPost, "/api/v1/apis", func(header http.Header, body []byte) (int, string) {
	// 	t.Logf("create api: %s\n", string(body))
	// 	return http.StatusOK, `{"code":"OK","api":{"id":"foo"}}`
	// })

	temp := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	res, err := client.CreateAPI(context.TODO(), kinde.CreateAPIParams{Name: temp, Audience: temp})
	assert.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.ID)

	t.Logf("created test api: %s\n", res.ID)

	res, err = client.GetAPI(context.TODO(), kinde.GetAPIParams{ID: res.ID})
	assert.NoError(t, err)
	require.NotNil(t, res)

	err = client.DeleteAPI(context.TODO(), kinde.DeleteAPIParams{ID: res.ID})
	assert.NoError(t, err)

}
