package clientcredentials

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/axatol/kinde-go/internal/logger"
)

// based on https://pkg.go.dev/golang.org/x/oauth2/clientcredentials
type OAuth2Transport struct {
	mu sync.RWMutex

	Domain       string
	Audience     string
	ClientID     string
	ClientSecret string
	Scope        []string
	Transport    http.RoundTripper
	Logger       logger.Logger
	Expiry       time.Time
	Token        string
}

type TokenExchangeResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (t *OAuth2Transport) RefreshToken(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	tokenEndpoint := t.Domain + "/oauth2/token"
	body := url.Values{
		"audience":      {t.Audience},
		"client_id":     {t.ClientID},
		"client_secret": {t.ClientSecret},
		"scope":         {strings.Join(t.Scope, " ")},
		"grant_type":    {"client_credentials"},
	}

	encodedBody := bytes.NewBuffer([]byte(body.Encode()))
	t.Logger.Logf("[OAuth2Transport.RefreshToken] %s - %s\n", http.MethodPost, tokenEndpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint, encodedBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	res, err := t.Transport.RoundTrip(req)
	if err != nil {
		return fmt.Errorf("failed to round trip request: %s", err)
	}

	t.Logger.Logf("[OAuth2Transport.RefreshToken] %s - %s - response status: %d\n", http.MethodPost, tokenEndpoint, res.StatusCode)

	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Logger.Logf("[OAuth2Transport.RefreshToken] %s - %s - response body: %s\n", http.MethodPost, tokenEndpoint, string(raw))
		return fmt.Errorf("unexpected status code: %d, body: %s", res.StatusCode, string(raw))
	}

	var response TokenExchangeResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return fmt.Errorf("failed to parse response body: %s", err)
	}

	if response.TokenType != "bearer" {
		return fmt.Errorf("unexpected token type: expected 'bearer', got '%s'", response.TokenType)
	}

	// leave 10 seconds to re-authenticate
	lifespan := time.Second * time.Duration(response.ExpiresIn-30)
	t.Expiry = time.Now().Add(lifespan)
	t.Token = response.AccessToken

	return nil
}

func (t *OAuth2Transport) GetToken(ctx context.Context) (string, error) {
	t.mu.RLock()
	if t.Transport == nil {
		t.Transport = http.DefaultTransport
	}

	if t.Logger == nil {
		t.Logger = logger.NoopLogger{}
	}

	expired := time.Now().After(t.Expiry)
	missing := t.Token == ""
	t.mu.RUnlock()

	if expired || missing {
		t.Logger.Logf("[Oauth2Transport.GetToken] refreshing token, expired: %v, missing: %v\n", expired, missing)
		if err := t.RefreshToken(ctx); err != nil {
			return "", fmt.Errorf("failed to retrieve access token: %s", err)
		}
	}

	return t.Token, nil
}

func (t *OAuth2Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	token, err := t.GetToken(r.Context())
	if err != nil {
		return nil, err
	}

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")

	return t.Transport.RoundTrip(r)
}
