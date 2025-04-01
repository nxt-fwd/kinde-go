package users

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/nxt-fwd/kinde-go/internal/client"
	"github.com/nxt-fwd/kinde-go/internal/phone"
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
		ID:          response.ID,
		FirstName:   response.GivenName,
		LastName:    response.FamilyName,
		IsSuspended: response.IsSuspended,
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

// addIdentityRequest handles the common API request logic for adding identities
func (c *Client) addIdentityRequest(ctx context.Context, userID string, params AddIdentityParams) (*Identity, error) {
	endpoint := fmt.Sprintf("/api/v1/users/%s/identities", userID)
	req, err := c.NewRequest(ctx, http.MethodPost, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response AddIdentityResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Identity, nil
}

// AddPhoneIdentity adds a phone identity to a user using a full international format phone number
func (c *Client) AddPhoneIdentity(ctx context.Context, userID string, fullPhoneNumber string) (*Identity, error) {
	localNumber, countryID, err := phone.ParseNumber(fullPhoneNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	return c.addIdentityRequest(ctx, userID, AddIdentityParams{
		Type:           IdentityTypePhone,
		Value:          localNumber,
		PhoneCountryID: countryID,
	})
}

// AddIdentity adds a new identity to a user
func (c *Client) AddIdentity(ctx context.Context, userID string, params AddIdentityParams) (*Identity, error) {
	if params.Type == IdentityTypePhone {
		// For phone identities, always convert to international format and use AddPhoneIdentity
		formattedNumber, err := phone.FormatNumber(params.Value, params.PhoneCountryID)
		if err != nil {
			return nil, fmt.Errorf("invalid phone number components: %w", err)
		}
		return c.AddPhoneIdentity(ctx, userID, formattedNumber)
	}

	return c.addIdentityRequest(ctx, userID, params)
}

// GetIdentities gets all identities for a user
func (c *Client) GetIdentities(ctx context.Context, userID string) ([]Identity, error) {
	endpoint := fmt.Sprintf("/api/v1/users/%s/identities", userID)
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var response struct {
		Code       string     `json:"code"`
		Message    string     `json:"message"`
		Identities []Identity `json:"identities"`
	}
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Identities, nil
}
