//go:build e2e
// +build e2e

package apis_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/axatol/kinde-go/api/apis"
	"github.com/axatol/kinde-go/api/applications"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2EList(t *testing.T) {
	client := apis.New(testutil.DefaultE2EClient(t))
	res, err := client.List(context.TODO())
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestE2EGet(t *testing.T) {
	client := apis.New(testutil.DefaultE2EClient(t))
	res, err := client.List(context.TODO())
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestE2ECreateGetAuthoriseDelete(t *testing.T) {
	client := apis.New(testutil.DefaultE2EClient(t))
	appClient := applications.New(testutil.DefaultE2EClient(t))
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	// Create a test application first
	app, err := appClient.Create(context.TODO(), applications.CreateParams{
		Name: tempID + "-app",
		Type: applications.TypeRegular,
	})
	assert.NoError(t, err)
	require.NotNil(t, app)
	require.NotEmpty(t, app.ID)
	t.Logf("created test application: %s\n", app.ID)

	// Ensure cleanup of test application
	defer func() {
		err := appClient.Delete(context.TODO(), app.ID)
		assert.NoError(t, err)
		t.Logf("deleted test application: %s\n", app.ID)
	}()

	// Create API
	res, err := client.Create(context.TODO(), apis.CreateParams{Name: tempID, Audience: tempID})
	assert.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.ID)

	id := res.ID
	t.Logf("created test api: %s\n", res.ID)

	res, err = client.Get(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, res)
	t.Logf("got test api: %+v\n", res)

	authoriseAppParams := apis.AuthorizeApplicationsParams{
		Applications: []apis.ApplicationAuthorization{{
			ID: app.ID,
		}},
	}

	err = client.AuthorizeApplications(context.TODO(), id, authoriseAppParams)
	assert.NoError(t, err)
	t.Logf("authorised test api application: %+v\n", authoriseAppParams)

	deauthoriseAppParams := apis.AuthorizeApplicationsParams{
		Applications: []apis.ApplicationAuthorization{{
			ID: app.ID, Operation: "delete",
		}},
	}

	err = client.AuthorizeApplications(context.TODO(), id, deauthoriseAppParams)
	assert.NoError(t, err)
	t.Logf("deauthorised test api application: %+v\n", deauthoriseAppParams)

	err = client.Delete(context.TODO(), id)
	assert.NoError(t, err)
	t.Logf("deleted test api: %+v\n", res)

	// Verify API deletion
	_, err = client.Get(context.TODO(), id)
	assert.Error(t, err)
	t.Logf("verified api deletion: %s\n", id)
}
