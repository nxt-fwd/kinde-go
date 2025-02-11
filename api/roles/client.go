package roles

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

// List roles
func (c *Client) List(ctx context.Context) ([]Role, error) {
	endpoint := "/api/v1/roles"
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response ListResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Roles, nil
}

// Create a new role
func (c *Client) Create(ctx context.Context, params CreateParams) (*Role, error) {
	endpoint := "/api/v1/roles"
	req, err := c.NewRequest(ctx, http.MethodPost, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response CreateResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Role, nil
}

// Get role details
func (c *Client) Get(ctx context.Context, id string) (*Role, error) {
	endpoint := fmt.Sprintf("/api/v1/roles/%s", id)
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Role    Role   `json:"role"`
	}
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	// Get role permissions
	perms, err := c.GetRolePermissions(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	response.Role.Permissions = perms

	return &response.Role, nil
}

// GetRolePermissions gets all permissions assigned to a role
func (c *Client) GetRolePermissions(ctx context.Context, roleID string) ([]string, error) {
	endpoint := fmt.Sprintf("/api/v1/roles/%s/permissions", roleID)
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Code        string       `json:"code"`
		Message     string       `json:"message"`
		NextToken   string       `json:"next_token"`
		Permissions []Permission `json:"permissions"`
	}
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	// Convert Permission objects to permission IDs
	var permissionIDs []string
	if response.Permissions != nil {
		for _, p := range response.Permissions {
			permissionIDs = append(permissionIDs, p.ID)
		}
	} else {
		permissionIDs = make([]string, 0)
	}

	return permissionIDs, nil
}

// Update role details
func (c *Client) Update(ctx context.Context, id string, params UpdateParams) (*Role, error) {
	endpoint := fmt.Sprintf("/api/v1/roles/%s", id)
	req, err := c.NewRequest(ctx, http.MethodPatch, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response UpdateResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Role, nil
}

// Delete a role
func (c *Client) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/api/v1/roles/%s", id)
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

// UpdatePermissions updates the permissions for a role
func (c *Client) UpdatePermissions(ctx context.Context, id string, params UpdatePermissionsParams) (*UpdatePermissionsResponse, error) {
	endpoint := fmt.Sprintf("/api/v1/roles/%s/permissions", id)
	req, err := c.NewRequest(ctx, http.MethodPatch, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response UpdatePermissionsResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// RemovePermission removes a specific permission from a role
func (c *Client) RemovePermission(ctx context.Context, roleID string, permissionID string) error {
	endpoint := fmt.Sprintf("/api/v1/roles/%s/permissions/%s", roleID, permissionID)
	req, err := c.NewRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}

	var response RemovePermissionResponse
	if err := c.DoRequest(req, &response); err != nil {
		return err
	}

	return nil
}

// ListPermissions lists all available permissions
func (c *Client) ListPermissions(ctx context.Context) ([]Permission, error) {
	endpoint := "/api/v1/permissions"
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response ListPermissionsResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Permissions, nil
} 
