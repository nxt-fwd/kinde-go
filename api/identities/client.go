package identities

import (
	"context"
	"fmt"
	"net/http"

	"github.com/nxt-fwd/kinde-go/api/users"
	"github.com/nxt-fwd/kinde-go/internal/client"
)

type Client struct {
	c client.Client
}

// UpdateIdentityRequest represents the request to update an identity
type UpdateIdentityRequest struct {
	IsPrimary bool `json:"is_primary"`
}

func New(client client.Client) *Client {
	return &Client{c: client}
}

// Get retrieves a specific identity by ID
func (c *Client) Get(ctx context.Context, identityID string) (*users.Identity, error) {
	endpoint := fmt.Sprintf("/api/v1/identities/%s", identityID)
	req, err := c.c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Code     string         `json:"code"`
		Message  string         `json:"message"`
		Identity users.Identity `json:"identity"`
	}
	if err := c.c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Identity, nil
}

// Update updates whether an identity is primary
func (c *Client) Update(ctx context.Context, identityID string, isPrimary bool) (*users.Identity, error) {
	endpoint := fmt.Sprintf("/api/v1/identities/%s", identityID)
	
	req := UpdateIdentityRequest{
		IsPrimary: isPrimary,
	}
	
	httpReq, err := c.c.NewRequest(ctx, http.MethodPatch, endpoint, nil, req)
	if err != nil {
		return nil, err
	}

	var response struct {
		Code     string         `json:"code"`
		Message  string         `json:"message"`
		Identity users.Identity `json:"identity"`
	}
	if err := c.c.DoRequest(httpReq, &response); err != nil {
		return nil, err
	}

	return &response.Identity, nil
}

// Delete deletes an identity
func (c *Client) Delete(ctx context.Context, identityID string) error {
	endpoint := fmt.Sprintf("/api/v1/identities/%s", identityID)
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