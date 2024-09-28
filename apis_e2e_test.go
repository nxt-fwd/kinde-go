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

func TestE2EListAPIs(t *testing.T) {
	client := defaultE2EClient(t)
	res, err := client.ListAPIs(context.TODO())
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestE2EGetAPI(t *testing.T) {
	client := defaultE2EClient(t)
	res, err := client.ListAPIs(context.TODO())
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestE2ECreateGetAuthoriseDeleteAPI(t *testing.T) {
	client := defaultE2EClient(t)
	tempID := fmt.Sprintf("test-%d", time.Now().UnixMilli())

	res, err := client.CreateAPI(context.TODO(), kinde.CreateAPIParams{Name: tempID, Audience: tempID})
	assert.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.ID)

	id := res.ID

	t.Logf("created test api: %s\n", res.ID)

	res, err = client.GetAPI(context.TODO(), id)
	assert.NoError(t, err)
	require.NotNil(t, res)

	t.Logf("got test api: %+v\n", res)

	authoriseAppParams := kinde.AuthorizeAPIApplicationsParams{
		Applications: []kinde.ApplicationAuthorization{{
			// Terraform Acceptance Example Application
			ID: "f61f05b791e142dcb44f113b54b2eee6",
		}},
	}

	err = client.AuthorizeAPIApplications(context.TODO(), id, authoriseAppParams)
	assert.NoError(t, err)

	t.Logf("authorised test api application: %+v\n", authoriseAppParams)

	deauthoriseAppParams := kinde.AuthorizeAPIApplicationsParams{
		Applications: []kinde.ApplicationAuthorization{{
			ID: "f61f05b791e142dcb44f113b54b2eee6", Operation: "delete",
		}},
	}

	err = client.AuthorizeAPIApplications(context.TODO(), id, deauthoriseAppParams)
	assert.NoError(t, err)

	t.Logf("deauthorised test api application: %+v\n", deauthoriseAppParams)

	err = client.DeleteAPI(context.TODO(), id)
	assert.NoError(t, err)

	t.Logf("deleted test api: %+v\n", res)
}
