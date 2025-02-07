package apis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nxt-fwd/kinde-go/internal/client"
)

type Client struct {
	client.Client
}

func New(client client.Client) *Client {
	return &Client{client}
}

type ListResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	NextToken string `json:"next_token"`
	APIs      []API  `json:"apis"`
}

// https://kinde.com/api/docs/#get-apis
//
// todo: pagination
func (c *Client) List(ctx context.Context) ([]API, error) {
	endpoint := "/api/v1/apis"
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response ListResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return response.APIs, nil
}

type CreateParams struct {
	Name     string `json:"name"`
	Audience string `json:"audience"`
}

type CreateResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	API     API    `json:"api"`
}

// https://kinde.com/api/docs/#create-api
//
// note: only ID will be populated
func (c *Client) Create(ctx context.Context, params CreateParams) (*API, error) {
	endpoint := "/api/v1/apis"
	req, err := c.NewRequest(ctx, http.MethodPost, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response CreateResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	api := API{ID: response.API.ID}
	return &api, nil
}

type GetResponse CreateResponse

// https://kinde.com/api/docs/#get-api
func (c *Client) Get(ctx context.Context, id string) (*API, error) {
	endpoint := fmt.Sprintf("/api/v1/apis/%s", id)
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response CreateResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.API, nil
}

type DeleteResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// https://kinde.com/api/docs/#delete-api
func (c *Client) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/api/v1/apis/%s", id)
	req, err := c.NewRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}

	var response CreateResponse
	if err := c.DoRequest(req, &response); err != nil {
		return err
	}

	return nil
}

type AuthorizeApplicationsParams struct {
	Applications []ApplicationAuthorization `json:"applications"`
}

type ApplicationAuthorization struct {
	ID string `json:"id"`
	// leave empty to assign, set to "delete" to unassign
	Operation string `json:"operation,omitempty"`
}

type AuthorizeApplicationsResponse struct {
	Code                     string   `json:"code"`
	Message                  string   `json:"message"`
	ApplicationsDisconnected []string `json:"applications_disconnected"`
	ApplicationsConnected    []string `json:"applications_connected"`
}

// https://kinde.com/api/docs/#authorize-api-applications
func (c *Client) AuthorizeApplications(ctx context.Context, id string, params AuthorizeApplicationsParams) error {
	endpoint := fmt.Sprintf("/api/v1/apis/%s/applications", id)
	req, err := c.NewRequest(ctx, http.MethodPatch, endpoint, nil, params)
	if err != nil {
		return err
	}

	var response AuthorizeApplicationsResponse
	if err := c.DoRequest(req, &response); err != nil {
		return err
	}

	return nil
}
