package permissions_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/nxt-fwd/kinde-go/api/permissions"
	"github.com/nxt-fwd/kinde-go/internal/client"
	"github.com/nxt-fwd/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := permissions.New(client.New(context.TODO(), nil))
	_, _ = client.List(context.TODO(), permissions.ListParams{})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/permissions"))
}

func TestCreate(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := permissions.New(client.New(context.TODO(), nil))
	_, _ = client.Create(context.TODO(), permissions.CreateParams{})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPost, "/api/v1/permissions"))
}

func TestUpdate(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := permissions.New(client.New(context.TODO(), nil))
	_ = client.Update(context.TODO(), "1", permissions.UpdateParams{})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPatch, "/api/v1/permissions/1"))
}

func TestDelete(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := permissions.New(client.New(context.TODO(), nil))
	_ = client.Delete(context.TODO(), "1")
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodDelete, "/api/v1/permissions/1"))
}
