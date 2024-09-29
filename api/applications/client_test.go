package applications_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/axatol/kinde-go/api/applications"
	"github.com/axatol/kinde-go/internal/client"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := applications.New(client.New(context.TODO(), nil))
	_, _ = client.List(context.TODO(), applications.ListParams{})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/applications"))
}

func TestCreate(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := applications.New(client.New(context.TODO(), nil))
	testServer.Handle(t, http.MethodPost, "/api/v1/applications", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, `{"name":"name","type":"reg"}`, string(body))
		return http.StatusOK, `{"code":"OK","application":{"id":"id","client_id":"client_id","client_secret":"client_secret"}}`
	})
	_, _ = client.Create(context.TODO(), applications.CreateParams{Name: "name", Type: applications.TypeRegular})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPost, "/api/v1/applications"))
}

func TestGet(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := applications.New(client.New(context.TODO(), nil))
	_, _ = client.Get(context.TODO(), "1")
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/applications/1"))
}

func TestDelete(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := applications.New(client.New(context.TODO(), nil))
	_ = client.Delete(context.TODO(), "1")
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodDelete, "/api/v1/applications/1"))
}
