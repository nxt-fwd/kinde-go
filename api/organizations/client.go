package organizations

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nxt-fwd/kinde-go/internal/client"
)

type Client struct {
	client.Client
}

func New(client client.Client) *Client {
	return &Client{client}
}

// List organizations
func (c *Client) List(ctx context.Context) ([]Organization, error) {
	endpoint := "/api/v1/organizations"
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response ListResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Organizations, nil
}

// Create a new organization
func (c *Client) Create(ctx context.Context, params CreateParams) (*Organization, error) {
	endpoint := "/api/v1/organization"
	req, err := c.NewRequest(ctx, http.MethodPost, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response CreateResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Organization, nil
}

// Get organization details
func (c *Client) Get(ctx context.Context, code string) (*Organization, error) {
	query := url.Values{}
	query.Set("code", code)

	endpoint := "/api/v1/organization"
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, query, nil)
	if err != nil {
		return nil, err
	}

	var organization Organization
	if err := c.DoRequest(req, &organization); err != nil {
		return nil, err
	}

	return &organization, nil
}

// Update organization details
func (c *Client) Update(ctx context.Context, code string, params UpdateParams) (*Organization, error) {
	endpoint := fmt.Sprintf("/api/v1/organization/%s", code)
	req, err := c.NewRequest(ctx, http.MethodPatch, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	// Get the updated organization
	return c.Get(ctx, code)
}

// Delete an organization
func (c *Client) Delete(ctx context.Context, code string) error {
	endpoint := fmt.Sprintf("/api/v1/organization/%s", code)
	req, err := c.NewRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}

	var response DeleteResponse
	if err := c.DoRequest(req, &response); err != nil {
		return err
	}

	return nil
}

// AddUsers adds users to an organization with specified roles and permissions
func (c *Client) AddUsers(ctx context.Context, code string, params AddUsersParams) error {
	endpoint := fmt.Sprintf("/api/v1/organizations/%s/users", code)
	req, err := c.NewRequest(ctx, http.MethodPost, endpoint, nil, params)
	if err != nil {
		return err
	}

	var response AddUsersResponse
	if err := c.DoRequest(req, &response); err != nil {
		return err
	}

	return nil
}

// AddUserRole adds a role to a user in an organization
func (c *Client) AddUserRole(ctx context.Context, orgCode string, userID string, roleID string) error {
	endpoint := fmt.Sprintf("/api/v1/organizations/%s/users/%s/roles", orgCode, userID)
	req, err := c.NewRequest(ctx, http.MethodPost, endpoint, nil, map[string]string{
		"role_id": roleID,
	})
	if err != nil {
		return err
	}

	var response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := c.DoRequest(req, &response); err != nil {
		return err
	}

	return nil
}

// GetUserRoles gets all roles assigned to a user in an organization
func (c *Client) GetUserRoles(ctx context.Context, orgCode string, userID string) ([]Role, error) {
	endpoint := fmt.Sprintf("/api/v1/organizations/%s/users/%s/roles", orgCode, userID)
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Roles   []Role `json:"roles"`
	}
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	// Initialize empty slice if no roles are returned
	if response.Roles == nil {
		response.Roles = make([]Role, 0)
	}

	return response.Roles, nil
}

// RemoveUserRole removes a role from a user in an organization
func (c *Client) RemoveUserRole(ctx context.Context, orgCode string, userID string, roleID string) error {
	endpoint := fmt.Sprintf("/api/v1/organizations/%s/users/%s/roles/%s", orgCode, userID, roleID)
	req, err := c.NewRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}

	var response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := c.DoRequest(req, &response); err != nil {
		return err
	}

	return nil
} 