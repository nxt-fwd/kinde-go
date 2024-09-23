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
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	res, err := client.CreateAPI(context.TODO(), kinde.CreateAPIParams{Name: tempID, Audience: tempID})
	assert.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.ID)

	t.Logf("created test api: %s\n", res.ID)

	res, err = client.GetAPI(context.TODO(), kinde.GetAPIParams{ID: res.ID})
	assert.NoError(t, err)
	require.NotNil(t, res)

	t.Logf("got test api: %+v\n", res)

	err = client.DeleteAPI(context.TODO(), kinde.DeleteAPIParams{ID: res.ID})
	assert.NoError(t, err)

}
