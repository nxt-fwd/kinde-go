package users_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/axatol/kinde-go/api/users"
	"github.com/axatol/kinde-go/internal/client"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := users.New(client.New(context.TODO(), nil))

	testServer.Handle(t, http.MethodGet, "/api/v1/users", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, "10", query.Get("page_size"))
		assert.Equal(t, "next_token_value", query.Get("next_token"))
		assert.Equal(t, "created_on", query.Get("sort"))
		return http.StatusOK, `{"code":"OK","users":[{"id":"test_id","first_name":"John","last_name":"Doe"}]}`
	})

	_, _ = client.List(context.TODO(), users.ListParams{
		PageSize:  10,
		NextToken: "next_token_value",
		Sort:      "created_on",
	})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/users"))
}

func TestCreate(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := users.New(client.New(context.TODO(), nil))

	testServer.Handle(t, http.MethodPost, "/api/v1/user", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, `{"profile":{"given_name":"John","family_name":"Doe","email":"john@example.com"}}`, string(body))
		return http.StatusOK, `{"id":"test_id","created":true,"identities":[]}`
	})

	testServer.Handle(t, http.MethodGet, "/api/v1/user", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, "test_id", query.Get("id"))
		return http.StatusOK, `{"id":"test_id","first_name":"John","last_name":"Doe"}`
	})

	_, _ = client.Create(context.TODO(), users.CreateParams{
		Profile: users.Profile{
			GivenName:  "John",
			FamilyName: "Doe",
			Email:      "john@example.com",
		},
	})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPost, "/api/v1/user"))
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/user"))
}

func TestGet(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := users.New(client.New(context.TODO(), nil))

	testServer.Handle(t, http.MethodGet, "/api/v1/user", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, "test_id", query.Get("id"))
		return http.StatusOK, `{"id":"test_id","first_name":"John","last_name":"Doe"}`
	})

	_, _ = client.Get(context.TODO(), "test_id")
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodGet, "/api/v1/user"))
}

func TestUpdate(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := users.New(client.New(context.TODO(), nil))

	testServer.Handle(t, http.MethodPatch, "/api/v1/user", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, "test_id", query.Get("id"))
		assert.Equal(t, `{"given_name":"John","family_name":"Doe"}`, string(body))
		return http.StatusOK, `{"id":"test_id","email":"john@example.com","picture":null,"given_name":"John","family_name":"Doe","is_suspended":false,"is_password_reset_requested":false}`
	})

	_, _ = client.Update(context.TODO(), "test_id", users.UpdateParams{
		GivenName:  "John",
		FamilyName: "Doe",
	})
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodPatch, "/api/v1/user"))
}

func TestDelete(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := users.New(client.New(context.TODO(), nil))

	testServer.Handle(t, http.MethodDelete, "/api/v1/user", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, "test_id", query.Get("id"))
		return http.StatusOK, `{"code":"OK","message":"User deleted"}`
	})

	_ = client.Delete(context.TODO(), "test_id")
	assert.Equal(t, 1, testServer.CallCount.Get(http.MethodDelete, "/api/v1/user"))
}

func TestAddToOrganization(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)
	client := users.New(client.New(context.TODO(), nil))

	testServer.Handle(t, http.MethodPost, "/api/v1/organizations/test_org/users", func(header http.Header, query url.Values, body []byte) (int, string) {
		assert.Equal(t, `{"users":[{"id":"test_id","roles":["manager"],"permissions":["admin"]}]}`, string(body))
		return http.StatusOK, `{"code":"OK","message":"Users successfully added to organization"}`
	})

	err := client.AddToOrganization(context.TODO(), "test_org", users.AddToOrgParams{
		Users: []users.AddToOrgUser{
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
