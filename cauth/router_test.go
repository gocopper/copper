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
	"github.com/stretchr/testify/assert"
)

func TestRouter_HandleSignup(t *testing.T) {
	t.Parallel()

	var (
		sessionResult cauth.SessionResult
		server        = httptest.NewServer(cauthtest.NewHandler(t))
	)

	defer server.Close()

	reqBody := strings.NewReader(`{
		"username": "test-user",
		"password": "test-pass"
	}`)

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		server.URL+"/api/auth/signup",
		reqBody,
	)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBodyJ, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	err = json.Unmarshal(respBodyJ, &sessionResult)
	assert.NoError(t, err)
}

func TestRouter_HandleLogin_Invalid(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(cauthtest.NewHandler(t))
	defer server.Close()

	reqBody := strings.NewReader(`{
		"username": "invalid-user",
		"password": "invalid-pass"
	}`)

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		server.URL+"/api/auth/login",
		reqBody,
	)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestRouter_HandleLogin(t *testing.T) {
	t.Parallel()

	var sessionResult cauth.SessionResult

	server := httptest.NewServer(cauthtest.NewHandler(t))
	defer server.Close()

	_ = cauthtest.CreateNewUserSession(t, server)

	reqBody := strings.NewReader(`{
		"username": "test-user",
		"password": "test-pass"
	}`)

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		server.URL+"/api/auth/login",
		reqBody,
	)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBodyJ, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	err = json.Unmarshal(respBodyJ, &sessionResult)
	assert.NoError(t, err)
}

func TestRouter_HandleLogout_InvalidSession(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(cauthtest.NewHandler(t))
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

	server := httptest.NewServer(cauthtest.NewHandler(t))
	defer server.Close()

	req, err := http.NewRequestWithContext(context.Background(),
		http.MethodPost,
		server.URL+"/api/auth/logout",
		nil,
	)
	assert.NoError(t, err)

	session := cauthtest.CreateNewUserSession(t, server)
	req.SetBasicAuth(session.Session.UUID, session.PlainSessionToken)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)

	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
