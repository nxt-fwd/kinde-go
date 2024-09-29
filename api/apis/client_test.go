package apis_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/axatol/kinde-go/api/apis"
	"github.com/axatol/kinde-go/internal/client"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := apis.New(client.New(context.TODO(), nil))
	_, _ = client.List(context.TODO())
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/apis"))
}

func TestCreate(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := apis.New(client.New(context.TODO(), nil))
	testServer.Handle(t, http.MethodPost, "/api/v1/apis", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, `{"name":"name","audience":"audience"}`, string(body))
		return http.StatusOK, `{"code":"OK","api":{"id":"id"}}`
	})
	_, _ = client.Create(context.TODO(), apis.CreateParams{Name: "name", Audience: "audience"})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPost, "/api/v1/apis"))
}

func TestGet(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := apis.New(client.New(context.TODO(), nil))
	_, _ = client.Get(context.TODO(), "1")
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/apis/1"))
}

func TestDelete(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := apis.New(client.New(context.TODO(), nil))
	_ = client.Delete(context.TODO(), "1")
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodDelete, "/api/v1/apis/1"))
}

func TestAuthorizeApplications(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := apis.New(client.New(context.TODO(), nil))
	_ = client.AuthorizeApplications(context.TODO(), "1", apis.AuthorizeApplicationsParams{})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPatch, "/api/v1/apis/1/applications"))
}
