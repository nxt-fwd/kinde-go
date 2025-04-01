package connections

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nxt-fwd/kinde-go/internal/client"
)

type Client struct {
	c client.Client
}

func New(client client.Client) *Client {
	return &Client{c: client}
}

// CreateParams represents the parameters for creating a new connection
type CreateParams struct {
	Name                string      `json:"name"`                           // Internal name of the connection
	DisplayName         string      `json:"display_name"`                   // Public facing name
	Strategy            Strategy    `json:"strategy"`                       // Identity provider identifier
	EnabledApplications []string    `json:"enabled_applications,omitempty"` // Client IDs of enabled applications
	Options             interface{} `json:"options,omitempty"`              // Connection-specific options
}

// CreateResponse represents the response from the create connection endpoint
type CreateResponse struct {
	Code       string     `json:"code"`
	Message    string     `json:"message"`
	Connection Connection `json:"connection"`
}

// ListResponse represents the response from the list connections endpoint
type ListResponse struct {
	Code        string       `json:"code"`
	Message     string       `json:"message"`
	HasMore     bool         `json:"has_more"`
	Connections []Connection `json:"connections"`
}

// GetResponse represents the response from the get connection endpoint
type GetResponse struct {
	Code       string     `json:"code"`
	Message    string     `json:"message"`
	Connection Connection `json:"connection"`
}

// Connection represents a connection in Kinde
type Connection struct {
	ID                  string      `json:"id"`
	Name                string      `json:"name"`
	DisplayName         string      `json:"display_name"`
	Strategy            string      `json:"strategy"`
	EnabledApplications []string    `json:"enabled_applications,omitempty"`
	Options             interface{} `json:"options,omitempty"`
}

// UpdateParams represents the parameters for updating a connection
type UpdateParams struct {
	Name                string      `json:"name,omitempty"`
	DisplayName         string      `json:"display_name,omitempty"`
	EnabledApplications []string    `json:"enabled_applications,omitempty"`
	Options             interface{} `json:"options,omitempty"`
}

// Create creates a new connection
func (c *Client) Create(ctx context.Context, params CreateParams) (*Connection, error) {
	endpoint := "/api/v1/connections"
	req, err := c.c.NewRequest(ctx, http.MethodPost, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response CreateResponse
	if err := c.c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Connection, nil
}

// List retrieves all connections
func (c *Client) List(ctx context.Context) ([]Connection, error) {
	endpoint := "/api/v1/connections"
	req, err := c.c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response ListResponse
	if err := c.c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Connections, nil
}

// Get retrieves a specific connection by ID
func (c *Client) Get(ctx context.Context, id string) (*Connection, error) {
	endpoint := fmt.Sprintf("/api/v1/connections/%s", id)
	req, err := c.c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response GetResponse
	if err := c.c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Connection, nil
}

// Update updates specific fields of a connection by ID.
//
// This endpoint uses PATCH to update only the specified fields of the connection.
// Fields that are not provided in the params will remain unchanged.
func (c *Client) Update(ctx context.Context, id string, params UpdateParams) (*Connection, error) {
	endpoint := fmt.Sprintf("/api/v1/connections/%s", id)
	req, err := c.c.NewRequest(ctx, http.MethodPatch, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := c.c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	// Get the updated connection
	return c.Get(ctx, id)
}

// ReplaceParams represents the parameters for replacing a connection
type ReplaceParams struct {
	Name                string      `json:"name"`                           // Internal name of the connection (required)
	DisplayName         string      `json:"display_name"`                   // Public facing name (required)
	EnabledApplications []string    `json:"enabled_applications,omitempty"` // Client IDs of enabled applications
	Options             interface{} `json:"options"`                        // Connection-specific options (required)
}

// Replace replaces all fields of a connection by ID.
//
// This endpoint uses PUT to replace the entire connection configuration.
// All required fields must be provided in the params, as this will overwrite
// the entire connection configuration.
func (c *Client) Replace(ctx context.Context, id string, params ReplaceParams) (*Connection, error) {
	endpoint := fmt.Sprintf("/api/v1/connections/%s", id)
	req, err := c.c.NewRequest(ctx, http.MethodPut, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := c.c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	// Get the updated connection
	return c.Get(ctx, id)
}

// Delete deletes a connection by ID
func (c *Client) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/api/v1/connections/%s", id)
	req, err := c.c.NewRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}

	var response struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := c.c.DoRequest(req, &response); err != nil {
		return err
	}

	return nil
}
