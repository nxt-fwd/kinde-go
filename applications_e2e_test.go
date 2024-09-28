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

func TestE2EListApplication(t *testing.T) {
	client := defaultE2EClient(t)
	res, err := client.ListApplications(context.TODO(), kinde.ListApplicationsParams{})
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestE2ECreateGetUpdateDeleteApplication(t *testing.T) {
	client := defaultE2EClient(t)
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	res, err := client.CreateApplication(context.TODO(), kinde.CreateApplicationParams{Name: tempID, Type: kinde.ApplicationTypeRegular})
	assert.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.ID)

	id := res.ID

	t.Logf("created test application: %s\n", id)

	res, err = client.GetApplication(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, res)

	t.Logf("got test application: %+v\n", res)

	updateParams := kinde.UpdateApplicationParams{
		Name:         tempID,
		LoginURI:     "https://example.com",
		HomepageURI:  "https://example.com",
		LogoutURIs:   []string{"https://example.com"},
		RedirectURIs: []string{"https://example.com"},
	}

	err = client.UpdateApplication(context.TODO(), id, updateParams)
	assert.NoError(t, err)

	t.Logf("updated test application: %+v\n", res)

	err = client.DeleteApplication(context.TODO(), id)
	assert.NoError(t, err)

	t.Logf("deleted test application: %+v\n", res)
}
