package testutil

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestServerConfig struct {
	Audience     string
	ClientID     string
	ClientSecret string
	GrantType    string
	AccessToken  string
	Scopes       []string
}

func DefaultTestServerConfig() TestServerConfig {
	return TestServerConfig{
		Audience:     "http://test",
		ClientID:     "123",
		ClientSecret: "456",
		GrantType:    "client_credentials",
		AccessToken:  "access_token",
		Scopes:       []string{},
	}
}

type TestServer struct {
	mux       *http.ServeMux
	Server    *httptest.Server
	Config    TestServerConfig
	CallCount CallCount
}

type TestServerHandler func(header http.Header, body []byte) (int, string)

func NewTestServer(t *testing.T, config TestServerConfig) *TestServer {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	helper := &TestServer{
		mux:    mux,
		Server: server,
		Config: config,
	}

	helper.Handle(t, http.MethodPost, "/oauth2/token", func(header http.Header, body []byte) (int, string) {
		payload, err := url.ParseQuery(string(body))
		assert.NoError(t, err)

		assert.Equal(t, config.Audience, payload.Get("audience"))
		assert.Equal(t, config.ClientID, payload.Get("client_id"))
		assert.Equal(t, config.ClientSecret, payload.Get("client_secret"))
		assert.Equal(t, config.GrantType, payload.Get("grant_type"))

		response := fmt.Sprintf(`{"access_token":"%s","token_type":"bearer"}`, config.AccessToken)
		return http.StatusOK, response
	})

	helper.HandleAuthenticated(t, "", "/", nil)

	return helper
}

func (s *TestServer) Handle(t *testing.T, method, path string, handler TestServerHandler) {
	s.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		callCount := s.CallCount.Get(r.Method, r.URL.Path)
		t.Logf("[TestServer.Handle] %s %s call count: %d -> %d\n", r.Method, r.URL.Path, callCount, callCount+1)
		s.CallCount.Inc(r.Method, r.URL.Path)

		if method != "" && r.Method != method {
			http.NotFound(w, r)
			return
		}

		if handler == nil {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"code":"OK"}`)
			return
		}

		defer r.Body.Close()
		raw, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		t.Logf("[TestServer.Handle] %s %s request body: %s\n", r.Method, r.URL.Path, string(raw))

		status, response := handler(r.Header, raw)
		t.Logf("[TestServer.Handle] %s %s response: %d - %s\n", r.Method, r.URL.Path, status, string(response))

		w.WriteHeader(status)
		fmt.Fprint(w, response)
	})
}

func (s *TestServer) HandleAuthenticated(t *testing.T, method, path string, handler TestServerHandler) {
	s.Handle(t, method, path, func(header http.Header, body []byte) (int, string) {
		expectedToken := fmt.Sprintf("Bearer %s", s.Config.AccessToken)
		actualToken := header.Get("authorization")
		assert.Equal(t, expectedToken, actualToken)
		if actualToken != expectedToken {
			return http.StatusUnauthorized, ""
		}

		if handler == nil {
			return http.StatusOK, `{"code":"OK"}`
		}

		return handler(header, body)
	})
}
