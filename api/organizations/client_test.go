package organizations_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/axatol/kinde-go/api/organizations"
	"github.com/axatol/kinde-go/internal/client"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := organizations.New(client.New(context.TODO(), nil))
	_, _ = client.List(context.TODO())
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/organizations"))
}

func TestCreate(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := organizations.New(client.New(context.TODO(), nil))
	testServer.Handle(t, http.MethodPost, "/api/v1/organization", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, `{"name":"name"}`, string(body))
		return http.StatusOK, `{"code":"OK","organization":{"code":"test_org","name":"name"}}`
	})
	_, _ = client.Create(context.TODO(), organizations.CreateParams{Name: "name"})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPost, "/api/v1/organization"))
}

func TestGet(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := organizations.New(client.New(context.TODO(), nil))
	testServer.Handle(t, http.MethodGet, "/api/v1/organization", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, "test_org", query.Get("code"))
		return http.StatusOK, `{"code":"OK","organization":{"code":"test_org","name":"Test Org"}}`
	})
	_, _ = client.Get(context.TODO(), "test_org")
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/organization"))
}

func TestUpdate(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := organizations.New(client.New(context.TODO(), nil))
	testServer.Handle(t, http.MethodPatch, "/api/v1/organization/test_org", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, `{"name":"updated","external_id":"test-id","background_color":"#fff","theme_code":"light"}`, string(body))
		return http.StatusOK, `{"code":"OK","organization":{"code":"test_org","name":"updated"}}`
	})
	_, _ = client.Update(context.TODO(), "test_org", organizations.UpdateParams{
		Name:            "updated",
		ExternalID:      "test-id",
		BackgroundColor: "#fff",
		ThemeCode:       "light",
	})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPatch, "/api/v1/organization/test_org"))
}

func TestDelete(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := organizations.New(client.New(context.TODO(), nil))
	testServer.Handle(t, http.MethodDelete, "/api/v1/organization/test_org", func(header http.Header, query url.Values, body []byte) (int, string) {
		return http.StatusOK, `{"code":"OK","message":"Organization deleted"}`
	})
	_ = client.Delete(context.TODO(), "test_org")
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodDelete, "/api/v1/organization/test_org"))
}

func TestAddUsers(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := organizations.New(client.New(context.TODO(), nil))

	testServer.Handle(t, http.MethodPost, "/api/v1/organizations/test_org/users", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, `{"users":[{"id":"test_id","roles":["manager"],"permissions":["admin"]}]}`, string(body))
		return http.StatusOK, `{"code":"OK","message":"Users successfully added to organization"}`
	})

	err := client.AddUsers(context.TODO(), "test_org", organizations.AddUsersParams{
		Users: []organizations.AddUser{
			{
				ID:          "test_id",
				Roles:       []string{"manager"},
				Permissions: []string{"admin"},
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPost, "/api/v1/organizations/test_org/users"))
} 