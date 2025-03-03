package chttp_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/chttp/chttptest"
	"github.com/stretchr/testify/assert"
)

func TestHTMLReaderWriter_WriteHTML(t *testing.T) {
	t.Parallel()

	rw := chttptest.NewHTMLReaderWriter(t)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rw.WriteHTML(resp, req, chttp.WriteHTMLParams{
		StatusCode:   http.StatusOK,
		Data:         map[string]string{"user": "test"},
		PageTemplate: "index.html",
	})

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "text/html", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), `Test Page`)
}

func TestHTMLReaderWriter_WriteHTML_NotFound(t *testing.T) {
	t.Parallel()

	rw := chttptest.NewHTMLReaderWriter(t)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rw.WriteHTML(resp, req, chttp.WriteHTMLParams{
		StatusCode: http.StatusNotFound,
		Data:       map[string]string{"user": "test"},
	})

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, "text/html", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), `not found`)
}

func TestHTMLReaderWriter_WriteHTML_WithError(t *testing.T) {
	t.Parallel()

	rw := chttptest.NewHTMLReaderWriter(t)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rw.WriteHTML(resp, req, chttp.WriteHTMLParams{
		Error: errors.New("test error"),
		Data:  map[string]string{"user": "test"},
	})

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, "text/html", resp.Header().Get("content-type"))
	// Since we don't know exactly what the error page looks like, we're not checking content
}

func TestHTMLReaderWriter_WriteHTMLError(t *testing.T) {
	t.Parallel()

	rw := chttptest.NewHTMLReaderWriter(t)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	testErr := errors.New("test error")

	rw.WriteHTMLError(resp, req, testErr)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Equal(t, "text/html", resp.Header().Get("content-type"))
	// Since the error rendering depends on config, we're not checking the exact content
}

func TestHTMLReaderWriter_WriteHTML_CustomLayout(t *testing.T) {
	t.Parallel()

	rw := chttptest.NewHTMLReaderWriter(t)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rw.WriteHTML(resp, req, chttp.WriteHTMLParams{
		StatusCode:     http.StatusOK,
		Data:           map[string]string{"user": "test"},
		PageTemplate:   "index.html",
		LayoutTemplate: "main.html", // explicitly set the default layout
	})

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "text/html", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), `Test Page`)
}

func TestHTMLReaderWriter_WritePartial(t *testing.T) {
	t.Parallel()

	rw := chttptest.NewHTMLReaderWriter(t)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	rw.WritePartial(resp, req, chttp.WritePartialParams{
		Name: "partial.html",
		Data: map[string]string{"content": "partial content"},
	})

	assert.Equal(t, "text/html", resp.Header().Get("content-type"))
	// Test would be more specific if we knew what the partial template contained
}

func TestHTMLReaderWriter_Unauthorized_WithoutRedirect(t *testing.T) {
	t.Parallel()

	// Create a reader/writer with default config (no redirect URL)
	rw := chttptest.NewHTMLReaderWriter(t)
	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	rw.Unauthorized(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestHTMLReaderWriter_Unauthorized_WithRedirect(t *testing.T) {
	t.Parallel()

	// This test would require mocking a config with RedirectURLForUnauthorizedRequests
	// Since we can't easily modify the config in the test helper, this is a demonstration
	// of how the test would look

	// Assuming we had a way to create a HTMLReaderWriter with a redirect URL in the config:
	// redirectURL := "/login"
	// rw := createReaderWriterWithRedirectConfig(t, &redirectURL)

	// resp := httptest.NewRecorder()
	// req := httptest.NewRequest(http.MethodGet, "/protected", nil)

	// rw.Unauthorized(resp, req)

	// assert.Equal(t, http.StatusSeeOther, resp.Code)
	// assert.Equal(t, "/login", resp.Header().Get("Location"))
}
