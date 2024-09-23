package kinde

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/axatol/kinde-go/internal/enum"
)

var (
	ErrApplicationTypeInvalid = errors.New("application type is invalid")
	ErrApplicationNotFound    = errors.New("application not found")
)

// https://kinde.com/api/docs/#kinde-management-api-applications
type Application struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Type         ApplicationType `json:"type"`
	ClientID     string          `json:"client_id"`
	ClientSecret string          `json:"client_secret"`
	LoginURI     string          `json:"login_uri"`
	HomepageURI  string          `json:"homepage_uri"`
}

var _ enum.Enum[ApplicationType] = (*ApplicationType)(nil)

type ApplicationType string

const (
	ApplicationTypeRegular               ApplicationType = "reg"
	ApplicationTypeSinglePageApplication ApplicationType = "spa"
	ApplicationTypeMachineToMachine      ApplicationType = "m2m"
)

func (t ApplicationType) Options() []ApplicationType {
	return []ApplicationType{
		ApplicationTypeRegular,
		ApplicationTypeSinglePageApplication,
		ApplicationTypeMachineToMachine,
	}
}

func (t ApplicationType) Valid(value string) error {
	return enum.Valid(t.Options(), value)
}

type ListApplicationsSortMethod string

const (
	ListApplicationsSortMethodNameAsc  ListApplicationsSortMethod = "name_asc"
	ListApplicationsSortMethodNameDesc ListApplicationsSortMethod = "name_desc"
)

func (t ListApplicationsSortMethod) Options() []ListApplicationsSortMethod {
	return []ListApplicationsSortMethod{
		ListApplicationsSortMethodNameAsc,
		ListApplicationsSortMethodNameDesc,
	}
}

func (t ListApplicationsSortMethod) Valid(value string) error {
	return enum.Valid(t.Options(), value)
}

type ListApplicationsParams struct {
	Sort      ListApplicationsSortMethod
	PageSize  int
	NextToken string
}

type ListApplicationsResponse struct {
	Code         string        `json:"code"`
	Message      string        `json:"message"`
	NextToken    string        `json:"next_token"`
	Applications []Application `json:"applications"`
}

// https://kinde.com/api/docs/#get-applications
//
// note: only id, name, and type will be populated
func (c *Client) ListApplications(ctx context.Context, params ListApplicationsParams) ([]Application, error) {
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

	var response ListApplicationsResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return response.Applications, nil
}

type CreateApplicationParams struct {
	Name string          `json:"name"`
	Type ApplicationType `json:"type"`
}

type CreateApplicationResponse struct {
	Code        string      `json:"code"`
	Message     string      `json:"message"`
	Application Application `json:"application"`
}

// https://kinde.com/api/docs/#create-application
//
// note: only id, client_id, and client_secret will be populated
func (c *Client) CreateApplication(ctx context.Context, params CreateApplicationParams) (*Application, error) {
	endpoint := "/api/v1/applications"
	req, err := c.NewRequest(ctx, http.MethodPost, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response CreateApplicationResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Application, nil
}

type GetApplicationParams struct {
	ID string
}

type GetApplicationResponse struct {
	Code        string      `json:"code"`
	Message     string      `json:"message"`
	Application Application `json:"application"`
}

// https://kinde.com/api/docs/#get-application
func (c *Client) GetApplication(ctx context.Context, params GetApplicationParams) (*Application, error) {
	endpoint := fmt.Sprintf("/api/v1/applications/%s", params.ID)
	req, err := c.NewRequest(ctx, http.MethodGet, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response GetApplicationResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	return &response.Application, nil
}

type UpdateApplicationParams struct {
	ID           string   `json:"-"`
	Name         string   `json:"name,omitempty"`
	LanguageKey  string   `json:"language_key,omitempty"`
	LogoutURIs   []string `json:"logout_uris,omitempty"`
	RedirectURIs []string `json:"redirect_uris,omitempty"`
	LoginURI     string   `json:"login_uri,omitempty"`
	HomepageURI  string   `json:"homepage_uri,omitempty"`
}

type UpdateApplicationResponse struct {
	ID           string
	Name         string   `json:"name"`
	LanguageKey  string   `json:"language_key"`
	LogoutURIs   []string `json:"logout_uris"`
	RedirectURIs []string `json:"redirect_uris"`
	LoginURI     string   `json:"login_uri"`
	HomepageURI  string   `json:"homepage_uri"`
}

func (c *Client) UpdateApplication(ctx context.Context, params UpdateApplicationParams) (*Application, error) {
	endpoint := fmt.Sprintf("/api/v1/applications/%s", params.ID)
	req, err := c.NewRequest(ctx, http.MethodPatch, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	var response UpdateApplicationResponse
	if err := c.DoRequest(req, &response); err != nil {
		return nil, err
	}

	// manually set ID because it's not returned in the response
	response.ID = params.ID

	return nil, nil
}

type DeleteApplicationParams struct {
	ID string
}

func (c *Client) DeleteApplication(ctx context.Context, params DeleteApplicationParams) error {
	endpoint := fmt.Sprintf("/api/v1/applications/%s", params.ID)
	req, err := c.NewRequest(ctx, http.MethodDelete, endpoint, nil, nil)
	if err != nil {
		return err
	}

	if err := c.DoRequest(req, nil); err != nil {
		return err
	}

	return nil
}
