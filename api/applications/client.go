package applications

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/axatol/kinde-go/internal/client"
	"github.com/axatol/kinde-go/internal/enum"
)

type Client struct {
	client.Client
}

func New(client client.Client) *Client {
	return &Client{client}
}

type ListParams struct {
	Sort      ListSortMethod
	PageSize  int
	NextToken string
}

type ListSortMethod string

const (
	ListSortMethodNameAsc  ListSortMethod = "name_asc"
	ListSortMethodNameDesc ListSortMethod = "name_desc"
)

func (t ListSortMethod) Options() []ListSortMethod {
	return []ListSortMethod{
		ListSortMethodNameAsc,
		ListSortMethodNameDesc,
	}
}

func (t ListSortMethod) Valid() error {
	return enum.Valid(t.Options(), t)
}

type ListResponse struct {
	Code         string        `json:"code"`
	Message      string        `json:"message"`
	NextToken    string        `json:"next_token"`
	Applications []Application `json:"applications"`
}

// https://kinde.com/api/docs/#get-applications
//
// note: only id, name, and type will be populated
func (c *Client) List(ctx context.Context, params ListParams) ([]Application, error) {
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

	endpoint := "/api/v1/applications"
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, query, nil)
	if err != nil {
		return nil, err
	}

	var response ListResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Applications, nil
}

type CreateParams struct {
	Name string `json:"name"`
	Type Type   `json:"type"`
}

type CreateResponse struct {
	Code        string      `json:"code"`
	Message     string      `json:"message"`
	Application Application `json:"application"`
}

// https://kinde.com/api/docs/#create-application
//
// note: client_secret will not be populated for spa applications
func (c *Client) Create(ctx context.Context, params CreateParams) (*Application, error) {
	endpoint := "/api/v1/applications"
	req, err := c.NewRequest(ctx, http.MethodPost, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response CreateResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Application, nil
}

type GetResponse struct {
	Code        string      `json:"code"`
	Message     string      `json:"message"`
	Application Application `json:"application"`
}

// https://kinde.com/api/docs/#get-application
func (c *Client) Get(ctx context.Context, id string) (*Application, error) {
	endpoint := fmt.Sprintf("/api/v1/applications/%s", id)
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response GetResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Application, nil
}

type UpdateParams struct {
	Name         string   `json:"name,omitempty"`
	LanguageKey  string   `json:"language_key,omitempty"`
	LogoutURIs   []string `json:"logout_uris,omitempty"`
	RedirectURIs []string `json:"redirect_uris,omitempty"`
	LoginURI     string   `json:"login_uri,omitempty"`
	HomepageURI  string   `json:"homepage_uri,omitempty"`
}

// https://kinde.com/api/docs/#update-application
//
// note: api doesn't return anything meaningful
func (c *Client) Update(ctx context.Context, id string, params UpdateParams) error {
	endpoint := fmt.Sprintf("/api/v1/applications/%s", id)
	req, err := c.NewRequest(ctx, http.MethodPatch, endpoint, nil, params)
	if err != nil {
		return err
	}

	if err := c.DoRequest(req, nil); err != nil {
		return err
	}

	return nil
}

// https://kinde.com/api/docs/#delete-application
func (c *Client) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/api/v1/applications/%s", id)
	req, err := c.NewRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}

	if err := c.DoRequest(req, nil); err != nil {
		return err
	}

	return nil
}
