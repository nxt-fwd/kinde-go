package kinde_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/axatol/kinde-go"
	"github.com/stretchr/testify/assert"
)

func TestGetAPIs(t *testing.T) {
	client, testServer := defaultScenario(t)
	_, _ = client.GetAPIs(context.TODO())
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/apis"))
}

func TestCreateAPI(t *testing.T) {
	client, testServer := defaultScenario(t)
	testServer.Handle(t, http.MethodPost, "/api/v1/apis", func(header http.Header, body []byte) (int, string) {
		assert.Equal(t, `{"name":"name","audience":"audience"}`, string(body))
		return http.StatusOK, `{"code":"OK","api":{"id":"id"}}`
	})
	_, _ = client.CreateAPI(context.TODO(), kinde.CreateAPIParams{Name: "name", Audience: "audience"})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPost, "/api/v1/apis"))
}

func TestGetAPI(t *testing.T) {
	client, testServer := defaultScenario(t)
	_, _ = client.GetAPI(context.TODO(), kinde.GetAPIParams{ID: "1"})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/apis/1"))
}

func TestDeleteAPI(t *testing.T) {
	client, testServer := defaultScenario(t)
	_ = client.DeleteAPI(context.TODO(), kinde.DeleteAPIParams{ID: "1"})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodDelete, "/api/v1/apis/1"))
}

func TestAuthorizeAPIApplications(t *testing.T) {
	client, testServer := defaultScenario(t)
	_ = client.AuthorizeAPIApplications(context.TODO(), kinde.AuthorizeAPIApplicationsParams{ID: "1"})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPatch, "/api/v1/apis/1/applications"))
}
