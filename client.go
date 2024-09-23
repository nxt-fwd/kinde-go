package kinde

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/axatol/kinde-go/internal/clientcredentials"
	"github.com/axatol/kinde-go/internal/logger"
	"github.com/tidwall/gjson"
)

type Client struct {
	client *http.Client
	domain string
	logger logger.Logger
}

func New(ctx context.Context, options *ClientOptions) *Client {
	if options == nil {
		options = NewClientOptions()
	}

	transport := &clientcredentials.OAuth2Transport{
		Domain:       options.Domain,
		Audience:     options.Audience,
		ClientID:     options.ClientID,
		ClientSecret: options.ClientSecret,
		Scope:        options.Scopes,
		Logger:       options.Logger,
	}

	client := &Client{
		client: &http.Client{Transport: transport},
		domain: options.Domain,
		logger: options.Logger,
	}

	return client
}

func (c *Client) NewRequest(ctx context.Context, method, path string, query url.Values, payload any) (*http.Request, error) {
	endpoint := strings.Builder{}
	endpoint.WriteString(c.domain)
	if !strings.HasPrefix(path, "/") {
		endpoint.WriteString("/")
	}
	endpoint.WriteString(path)

	if query != nil {
		endpoint.WriteString("?" + query.Encode())
	}

	c.logger.Logf("[Client.NewRequest] %s %s - url: %s\n", method, path, endpoint.String())

	var encoded io.Reader
	if payload != nil {
		raw, err := json.Marshal(payload)
		if err != nil {
			return nil, RequestError{
				Method: method,
				Path:   path,
				Err:    fmt.Errorf("failed to marshal request payload: %w", err),
			}
		}

		c.logger.Logf("[Client.NewRequest] %s %s - request payload: %s\n", method, path, string(raw))
		encoded = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), encoded)
	if err != nil {
		return nil, RequestError{
			Method: method,
			Path:   path,
			Err:    fmt.Errorf("failed to create request: %w", err),
		}
	}

	c.logger.Logf("[Client.NewRequest] %s %s - request payload content length: %d\n", method, path, req.ContentLength)

	return req, nil
}

func (c *Client) DoRequest(req *http.Request, result any) error {
	res, err := c.client.Do(req)
	if err != nil {
		return RequestError{
			Method:     req.Method,
			Path:       req.URL.Path,
			StatusCode: res.StatusCode,
			Err:        fmt.Errorf("failed to execute request: %w", err),
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
