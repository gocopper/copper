package cauth_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gocopper/copper/cauth"
	"github.com/gocopper/copper/cauth/cauthtest"
	"github.com/gocopper/copper/chttp"
	"github.com/stretchr/testify/assert"
)

func TestRouter_HandleSignup(t *testing.T) {
	t.Parallel()

	var (
		sessionResult cauth.SessionResult
		router        = cauthtest.NewRouter(t)
	)

	reqBody := strings.NewReader(`{
		"username": "test-user",
		"password": "test-pass"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", reqBody)
	resp := httptest.NewRecorder()

	http.HandlerFunc(router.HandleSignup).ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	respBodyJ, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	err = json.Unmarshal(respBodyJ, &sessionResult)
	assert.NoError(t, err)
}

func TestRouter_HandleLogin_Invalid(t *testing.T) {
	t.Parallel()

	router := cauthtest.NewRouter(t)

	reqBody := strings.NewReader(`{
		"username": "invalid-user",
		"password": "invalid-pass"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", reqBody)
	resp := httptest.NewRecorder()

	http.HandlerFunc(router.HandleLogin).ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestRouter_HandleLogin(t *testing.T) {
	t.Parallel()

	var sessionResult cauth.SessionResult

	router := cauthtest.NewRouter(t)

	_ = cauthtest.CreateNewUserSession(t, router)

	reqBody := strings.NewReader(`{
		"username": "test-user",
		"password": "test-pass"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", reqBody)
	resp := httptest.NewRecorder()

	http.HandlerFunc(router.HandleLogin).ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	respBodyJ, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	err = json.Unmarshal(respBodyJ, &sessionResult)
	assert.NoError(t, err)
}

func TestRouter_HandleLogout_InvalidSession(t *testing.T) {
	t.Parallel()

	router := cauthtest.NewRouter(t)

	server := httptest.NewServer(chttp.NewHandler(chttp.NewHandlerParams{
		Routers:           []chttp.Router{router},
		GlobalMiddlewares: nil,
	}))
	defer server.Close()

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		server.URL+"/api/auth/logout",
		nil,
	)
	assert.NoError(t, err)

	req.SetBasicAuth("invalid-uuid", "invalid-token")

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestRouter_HandleLogout(t *testing.T) {
	t.Parallel()

	router := cauthtest.NewRouter(t)
	session := cauthtest.CreateNewUserSession(t, router)

	server := httptest.NewServer(chttp.NewHandler(chttp.NewHandlerParams{
		Routers:           []chttp.Router{router},
		GlobalMiddlewares: nil,
	}))
	defer server.Close()

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		server.URL+"/api/auth/logout",
		nil,
	)
	assert.NoError(t, err)

	req.SetBasicAuth(session.Session.UUID, session.PlainSessionToken)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
