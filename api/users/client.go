package users

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/axatol/kinde-go/internal/client"
)

type Client struct {
	client.Client
}

func New(client client.Client) *Client {
	return &Client{client}
}

// List users
func (c *Client) List(ctx context.Context, params ListParams) ([]User, error) {
	query := url.Values{}
	if params.PageSize > 0 {
		query.Set("page_size", fmt.Sprintf("%d", params.PageSize))
	}
	if params.NextToken != "" {
		query.Set("next_token", params.NextToken)
	}
	if params.Sort != "" {
		query.Set("sort", params.Sort)
	}

	endpoint := "/api/v1/users"
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, query, nil)
	if err != nil {
		return nil, err
	}

	var response ListResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Users, nil
}

// Create a new user
func (c *Client) Create(ctx context.Context, params CreateParams) (*User, error) {
	endpoint := "/api/v1/user"
	req, err := c.NewRequest(ctx, http.MethodPost, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response CreateResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	// After successful creation, get the user details
	return c.Get(ctx, response.ID)
}

// Get user details
func (c *Client) Get(ctx context.Context, id string) (*User, error) {
	query := url.Values{}
	query.Set("id", id)

	endpoint := "/api/v1/user"
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, query, nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := c.DoRequest(req, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// Update user details
func (c *Client) Update(ctx context.Context, id string, params UpdateParams) (*User, error) {
	query := url.Values{}
	query.Set("id", id)

	endpoint := "/api/v1/user"
	req, err := c.NewRequest(ctx, http.MethodPatch, endpoint, query, params)
	if err != nil {
		return nil, err
	}

	var response UpdateResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	// Convert UpdateResponse to User
	user := &User{
		ID:           response.ID,
		FirstName:    response.GivenName,
		LastName:     response.FamilyName,
		IsSuspended:  response.IsSuspended,
	}
	if response.Email != nil {
		user.PreferredEmail = *response.Email
	}
	if response.Picture != nil {
		user.Picture = *response.Picture
	}

	return user, nil
}

// Delete a user
func (c *Client) Delete(ctx context.Context, id string) error {
	query := url.Values{}
	query.Set("id", id)

	endpoint := "/api/v1/user"
	req, err := c.NewRequest(ctx, http.MethodDelete, endpoint, query, nil)
	if err != nil {
		return err
	}

	var response DeleteResponse
	if err := c.DoRequest(req, &response); err != nil {
		return err
	}

	return nil
} 
