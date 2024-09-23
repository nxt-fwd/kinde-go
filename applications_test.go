package kinde_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/axatol/kinde-go"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestListApplications(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := kinde.New(context.TODO(), nil)
	_, _ = client.ListApplications(context.TODO(), kinde.ListApplicationsParams{})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/applications"))
}

func TestCreateApplication(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := kinde.New(context.TODO(), nil)
	testServer.Handle(t, http.MethodPost, "/api/v1/applications", func(header http.Header, body []byte) (int, string) {
		assert.Equal(t, `{"name":"name","type":"reg"}`, string(body))
		return http.StatusOK, `{"code":"OK","application":{"id":"id","client_id":"client_id","client_secret":"client_secret"}}`
	})
	_, _ = client.CreateApplication(context.TODO(), kinde.CreateApplicationParams{Name: "name", Type: kinde.ApplicationTypeRegular})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPost, "/api/v1/applications"))
}

func TestGetApplication(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := kinde.New(context.TODO(), nil)
	_, _ = client.GetApplication(context.TODO(), kinde.GetApplicationParams{ID: "1"})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/applications/1"))
}

func TestDeleteApplication(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := kinde.New(context.TODO(), nil)
	_ = client.DeleteApplication(context.TODO(), kinde.DeleteApplicationParams{ID: "1"})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodDelete, "/api/v1/applications/1"))
}
