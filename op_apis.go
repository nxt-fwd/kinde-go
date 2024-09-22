package kinde

import (
	"context"
	"fmt"
	"net/http"
)

// https://kinde.com/api/docs/#kinde-management-api-apis
type API struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Audience        string        `json:"audience"`
	IsManagementApi bool          `json:"is_management_api"`
	Applications    []Application `json:"applications,omitempty"`
}

type ListAPIResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	NextToken string `json:"next_token"`
	APIs      []API  `json:"apis"`
}

// https://kinde.com/api/docs/#get-apis
//
// todo: pagination
func (c *Client) GetAPIs(ctx context.Context) ([]API, error) {
	endpoint := "/api/v1/apis"
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response ListAPIResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return response.APIs, nil
}

type CreateAPIParams struct {
	Name     string `json:"name"`
	Audience string `json:"audience"`
}

type CreateAPIResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	API     API    `json:"api"`
}

// https://kinde.com/api/docs/#create-api
//
// note: only ID will be populated
func (c *Client) CreateAPI(ctx context.Context, params CreateAPIParams) (*API, error) {
	endpoint := "/api/v1/apis"
	req, err := c.NewRequest(ctx, http.MethodPost, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response CreateAPIResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	api := API{ID: response.API.ID}
	return &api, nil
}

type GetAPIParams struct {
	ID string
}

type GetAPIResponse CreateAPIResponse

// https://kinde.com/api/docs/#get-api
func (c *Client) GetAPI(ctx context.Context, params GetAPIParams) (*API, error) {
	endpoint := fmt.Sprintf("/api/v1/apis/%s", params.ID)
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response CreateAPIResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.API, nil
}

type DeleteAPIParams struct {
	ID string
}

type DeleteAPIResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// https://kinde.com/api/docs/#delete-api
func (c *Client) DeleteAPI(ctx context.Context, params DeleteAPIParams) error {
	endpoint := fmt.Sprintf("/api/v1/apis/%s", params.ID)
	req, err := c.NewRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}

	var response CreateAPIResponse
	if err := c.DoRequest(req, &response); err != nil {
		return err
	}

	return nil
}

type AuthorizeAPIApplicationsParams struct {
	ID           string                     `json:"-"`
	Applications []ApplicationAuthorization `json:"applications"`
}

type ApplicationAuthorization struct {
	ID string `json:"id"`
	// leave empty to assign, set to "delete" to unassign
	Operation string `json:"operation,omitempty"`
}

type AuthorizeAPIApplicationsResponse struct {
	Code                     string   `json:"code"`
	Message                  string   `json:"message"`
	ApplicationsDisconnected []string `json:"applications_disconnected"`
	ApplicationsConnected    []string `json:"applications_connected"`
}

// https://kinde.com/api/docs/#authorize-api-applications
func (c *Client) AuthorizeAPIApplications(ctx context.Context, params AuthorizeAPIApplicationsParams) error {
	endpoint := fmt.Sprintf("/api/v1/apis/%s/applications", params.ID)
	req, err := c.NewRequest(ctx, http.MethodPatch, endpoint, nil, params)
	if err != nil {
		return err
	}

	var response AuthorizeAPIApplicationsResponse
	if err := c.DoRequest(req, &response); err != nil {
		return err
	}

	return nil
}
