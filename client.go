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

type KindeErrors []KindeError

func (errs KindeErrors) Error() string {
	messages := make([]string, 0, len(errs))
	for _, err := range errs {
		messages = append(messages, err.Error())
	}

	return strings.Join(messages, ", ")
}

func (errs KindeErrors) Has(code string) bool {
	for _, err := range errs {
		if err.Code == code {
			return true
		}
	}

	return false
}

func (errs *KindeErrors) UnmarshalJSON(data []byte) error {
	rawErrs := gjson.GetBytes(data, "errors")

	// its possible the kinde api may return an array or a single object
	// we need to handle both

	if rawErrs.IsArray() {
		var target struct {
			Errors []KindeError `json:"errors"`
		}

		if err := json.Unmarshal(data, &target); err != nil {
			return fmt.Errorf("failed to parse error list: %s", err)
		}

		*errs = target.Errors
		return nil
	}

	var target struct {
		Errors KindeError `json:"errors"`
	}

	if err := json.Unmarshal(data, &target); err != nil {
		return fmt.Errorf("failed to parse error: %s", err)
	}

	*errs = KindeErrors{target.Errors}
	return nil
}

type KindeError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e KindeError) Error() string {
	builder := strings.Builder{}
	builder.WriteString(e.Code)
	builder.WriteString(": ")
	if e.Message != "" {
		builder.WriteString(e.Message)
	} else {
		builder.WriteString("N/A")
	}

	return builder.String()
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
			return nil, fmt.Errorf("failed to marshal request payload: %s", err)
		}

		c.logger.Logf("[Client.NewRequest] %s %s - request payload: %s\n", method, path, string(raw))
		encoded = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	c.logger.Logf("[Client.NewRequest] %s %s - request payload content length: %d\n", method, path, req.ContentLength)

	return req, nil
}

func (c *Client) DoRequest(req *http.Request, result any) error {
	res, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %s", err)
	}

	c.logger.Logf("[Client.DoRequest] %s %s - response status: %d\n", req.Method, req.URL.Path, res.StatusCode)

	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	c.logger.Logf("[Client.DoRequest] %s %s - response body: %s\n", req.Method, req.URL.Path, string(raw))

	var errs KindeErrors
	if rawErrs := gjson.GetBytes(raw, "errors"); rawErrs.Exists() {
		if err := json.Unmarshal(raw, &errs); err != nil {
			return fmt.Errorf("failed to parse error response body: %s", err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("request failed: %s", errs)
	}

	// probably wont happen but just in case
	if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unexpected status code %d: %s", res.StatusCode, raw)
	}

	if err := json.Unmarshal(raw, result); err != nil {
		return fmt.Errorf("failed to parse response body: %s", err)
	}

	return nil
}
