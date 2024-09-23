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

func TestE2ECreateListDeleteApplication(t *testing.T) {
	client := defaultE2EClient(t)
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	res, err := client.CreateApplication(context.TODO(), kinde.CreateApplicationParams{Name: tempID, Type: kinde.ApplicationTypeRegular})
	assert.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.ID)

	t.Logf("created test application: %s\n", res.ID)

	res, err = client.GetApplication(context.TODO(), kinde.GetApplicationParams{ID: res.ID})
	assert.NoError(t, err)
	require.NotNil(t, res)

	t.Logf("got test application: %+v\n", res)

	err = client.DeleteApplication(context.TODO(), kinde.DeleteApplicationParams{ID: res.ID})
	assert.NoError(t, err)
}
