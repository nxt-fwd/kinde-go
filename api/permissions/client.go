package permissions

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/axatol/kinde-go/internal/client"
	"github.com/axatol/kinde-go/internal/enum"
)

var (
	ErrPermissionNotFound = fmt.Errorf("permission not found")
)

type Client struct {
	client.Client
}

func New(client client.Client) *Client {
	return &Client{Client: client}
}

type ListParams struct {
	Sort      ListSortMethod
	PageSize  int
	NextToken string
}

var _ enum.Enum[ListSortMethod] = (*ListSortMethod)(nil)

type ListSortMethod string

const (
	ListSortNameAsc  ListSortMethod = "name_asc"
	ListSortNameDesc ListSortMethod = "name_desc"
	ListSortIDAsc    ListSortMethod = "id_asc"
	ListSortIDDesc   ListSortMethod = "id_desc"
)

func (t ListSortMethod) Options() []ListSortMethod {
	return []ListSortMethod{
		ListSortNameAsc,
		ListSortNameDesc,
		ListSortIDAsc,
		ListSortIDDesc,
	}
}

func (t ListSortMethod) Valid() error {
	return enum.Valid(t.Options(), t)
}

type ListResponse struct {
	Code        string       `json:"code"`
	Message     string       `json:"message"`
	NextToken   string       `json:"next_token"`
	Permissions []Permission `json:"permissions"`
}

func (r ListResponse) GetNextToken() string { return r.NextToken }

func (r ListResponse) GetData() []Permission { return r.Permissions }

// https://kinde.com/api/docs/#list-permissions
func (c *Client) List(ctx context.Context, params ListParams) ([]Permission, error) {
	query := url.Values{}
	if params.Sort != "" {
		query.Set("sort", string(params.Sort))
	}

	if params.PageSize > 0 {
		query.Set("page_size", fmt.Sprint(params.PageSize))
	}

	if params.NextToken != "" {
		query.Set("next_token", params.NextToken)
	}

	endpoint := "/api/v1/permissions"
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, query, nil)
	if err != nil {
		return nil, err
	}

	var response ListResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Permissions, nil
}

type SearchParams struct {
	Name string
	Key  string
}

func (c *Client) Search(ctx context.Context, params SearchParams) (*Permission, error) {
	opts := client.PaginatorOptions{
		PageSize: 100,
		Sort:     string(ListSortNameAsc),
	}

	paginator := client.NewPaginator[Permission, ListResponse](c, "/api/v1/permissions", opts)

	for paginator.HasNext() {
		permissions, err := paginator.Next(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to find permission with name %s and key %s: %w", params.Name, params.Key, err)
		}

		for _, permission := range permissions {
			if permission.Name == params.Name && permission.Key == params.Key {
				return &permission, nil
			}
		}
	}

	return nil, fmt.Errorf("could not find permission with name %s and key %s: %w", params.Name, params.Key, ErrPermissionNotFound)
}

type CreateParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Key         string `json:"key"`
}

type CreateResponse struct {
	Code       string     `json:"code"`
	Message    string     `json:"message"`
	Permission Permission `json:"permission"`
}

// https://kinde.com/api/docs/#create-permission
//
// note: only ID will be populated
func (c *Client) Create(ctx context.Context, params CreateParams) (*Permission, error) {
	endpoint := "/api/v1/permissions"
	req, err := c.NewRequest(ctx, http.MethodPost, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response CreateResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Permission, nil
}

type UpdateParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Key         string `json:"key"`
}

// https://kinde.com/api/docs/#update-permission
func (c *Client) Update(ctx context.Context, id string, params UpdateParams) error {
	endpoint := fmt.Sprintf("/api/v1/permissions/%s", id)
	req, err := c.NewRequest(ctx, http.MethodPatch, endpoint, nil, params)
	if err != nil {
		return err
	}

	if err := c.DoRequest(req, nil); err != nil {
		return err
	}

	return nil
}

// https://kinde.com/api/docs/#delete-permission
func (c *Client) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/api/v1/permissions/%s", id)
	req, err := c.NewRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}

	if err := c.DoRequest(req, nil); err != nil {
		return err
	}

	return nil
}
