package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nxt-fwd/kinde-go/internal/logger"
	"github.com/nxt-fwd/kinde-go/internal/oauth2"
	"github.com/tidwall/gjson"
)

type Client interface {
	NewRequest(ctx context.Context, method, path string, query url.Values, payload any) (*http.Request, error)
	DoRequest(req *http.Request, result any) error
}

type clientImpl struct {
	client  *http.Client
	domain  string
	options *ClientOptions
	logger  logger.Logger
}

func New(ctx context.Context, options *ClientOptions) Client {
	if options == nil {
		options = NewClientOptions()
	}

	if err := options.Validate(); err != nil {
		// Return a client that will always return the validation error
		return &errorClient{err: err}
	}

	transport := &oauth2.OAuth2Transport{
		Domain:       options.Domain,
		Audience:     options.Audience,
		ClientID:     options.ClientID,
		ClientSecret: options.ClientSecret,
		Scope:        options.Scopes,
		Logger:       options.Logger,
	}

	client := &clientImpl{
		client: &http.Client{Transport: transport},
		domain: options.Domain,
		options: options,
		logger: options.Logger,
	}

	return client
}

func (c *clientImpl) NewRequest(ctx context.Context, method string, path string, query url.Values, body any) (*http.Request, error) {
	// Validate credentials before proceeding
	if err := c.options.Validate(); err != nil {
		return nil, fmt.Errorf("invalid client configuration: %w", err)
	}

	// Build URL
	u, err := url.Parse(c.domain)
	if err != nil {
		return nil, fmt.Errorf("failed to parse domain URL: %w", err)
	}

	u.Path = path
	if query != nil {
		u.RawQuery = query.Encode()
	}

	// Create request
	var buf io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		buf = bytes.NewBuffer(raw)
		c.logger.Logf("[Client.NewRequest] %s %s - request body: %s\n", method, path, string(raw))
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if token := c.options.GetAccessToken(); token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *clientImpl) DoRequest(req *http.Request, result any) error {
	res, err := c.client.Do(req)
	if err != nil {
		return RequestError{
			Method:     req.Method,
			Path:       req.URL.Path,
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("failed to execute request: %w", err),
		}
	}

	if res == nil {
		return RequestError{
			Method:     req.Method,
			Path:       req.URL.Path,
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("received nil response from server"),
		}
	}

	c.logger.Logf("[Client.DoRequest] %s %s - response status: %d\n", req.Method, req.URL.Path, res.StatusCode)

	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return RequestError{
			Method:     req.Method,
			Path:       req.URL.Path,
			StatusCode: res.StatusCode,
			Err:        fmt.Errorf("failed to read response body: %w", err),
		}
	}

	c.logger.Logf("[Client.DoRequest] %s %s - response body: %s\n", req.Method, req.URL.Path, string(raw))

	// Handle authentication errors specifically
	if res.StatusCode == http.StatusUnauthorized {
		return RequestError{
			Method:     req.Method,
			Path:       req.URL.Path,
			StatusCode: res.StatusCode,
			Err:        fmt.Errorf("authentication failed: invalid credentials or token"),
		}
	}

	var errs KindeErrors
	if rawErrs := gjson.GetBytes(raw, "errors"); rawErrs.Exists() {
		if err := json.Unmarshal(raw, &errs); err != nil {
			return RequestError{
				Method:     req.Method,
				Path:       req.URL.Path,
				StatusCode: res.StatusCode,
				Err:        fmt.Errorf("failed to parse error response body: %w", err),
			}
		}
	}

	if len(errs) > 0 {
		return RequestError{
			Method:     req.Method,
			Path:       req.URL.Path,
			StatusCode: res.StatusCode,
			Err:        fmt.Errorf("request failed: %w", errs),
		}
	}

	// probably wont happen but just in case
	if res.StatusCode >= http.StatusBadRequest {
		return RequestError{
			Method:     req.Method,
			Path:       req.URL.Path,
			StatusCode: res.StatusCode,
			Err:        fmt.Errorf("unexpected status code %d: %s", res.StatusCode, string(raw)),
		}
	}

	if result != nil {
		if err := json.Unmarshal(raw, result); err != nil {
			return RequestError{
				Method:     req.Method,
				Path:       req.URL.Path,
				StatusCode: res.StatusCode,
				Err:        fmt.Errorf("failed to parse response body: %w", err),
			}
		}
	}

	return nil
}

// errorClient is a client implementation that always returns the same error
type errorClient struct {
	err error
}

func (c *errorClient) NewRequest(ctx context.Context, method, path string, query url.Values, payload any) (*http.Request, error) {
	return nil, c.err
}

func (c *errorClient) DoRequest(req *http.Request, result any) error {
	return c.err
}
