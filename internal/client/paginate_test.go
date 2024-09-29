package client_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/axatol/kinde-go/internal/client"
	"github.com/axatol/kinde-go/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testPaginationResponse struct {
	NextToken string `json:"next_token"`
	Data      []any  `json:"data"`
}

func (r testPaginationResponse) GetNextToken() string { return r.NextToken }

func (r testPaginationResponse) GetData() []any { return r.Data }

func TestPagination(t *testing.T) {
	testServer := testutil.NewTestServer(t, nil)

	options := client.PaginatorOptions{
		Sort:     "name_asc",
		PageSize: 10,
	}

	testServer.HandleAuthenticated(t, http.MethodGet, "/api/v1/pagination", func(header http.Header, query url.Values, body []byte) (int, string) {
		if testServer.CallCount.Get(http.MethodGet, "/api/v1/pagination") == 1 {
			assert.Equal(t, "10", query.Get("page_size"))
			assert.Equal(t, "name_asc", query.Get("sort"))
			assert.Equal(t, "", query.Get("next_token"))
			return http.StatusOK, `{"code":"OK","data":[{"id":"1"}],"next_token":"next_token"}`
		}

		if testServer.CallCount.Get(http.MethodGet, "/api/v1/pagination") == 2 {
			assert.Equal(t, "10", query.Get("page_size"))
			assert.Equal(t, "name_asc", query.Get("sort"))
			assert.Equal(t, "next_token", query.Get("next_token"))
			return http.StatusOK, `{"code":"OK","data":null,"next_token":null}`
		}

		require.FailNow(t, "unexpected call")
		return 0, ""
	})

	paginator := client.NewPaginator[any, testPaginationResponse](client.New(context.TODO(), nil), "/api/v1/pagination", options)
	assert.True(t, paginator.HasNext())
	data, err := paginator.Next(context.TODO())
	assert.NoError(t, err)
	require.NotNil(t, data)
	assert.Len(t, data, 1)

	assert.True(t, paginator.HasNext())
	data, err = paginator.Next(context.TODO())
	assert.NoError(t, err)
	require.NotNil(t, data)
	assert.Len(t, data, 0)

	assert.False(t, paginator.HasNext())
	assert.Equal(t, 2, testServer.CallCount.Get(http.MethodGet, "/api/v1/pagination"))
}
