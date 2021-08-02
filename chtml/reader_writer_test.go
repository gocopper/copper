package chtml_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocopper/copper/chtml"
	"github.com/gocopper/copper/chtml/chtmltest"

	"github.com/stretchr/testify/assert"
)

func TestReaderWriter_WriteHTML(t *testing.T) {
	t.Parallel()

	rw := chtmltest.NewReaderWriter(t)
	resp := httptest.NewRecorder()

	rw.WriteHTML(resp, httptest.NewRequest(http.MethodGet, "/", nil), chtml.WriteHTMLParams{
		StatusCode:   http.StatusOK,
		Data:         map[string]string{"user": "test"},
		PageTemplate: "index.html",
	})

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "text/html", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), `Test Page`)
}

func TestReaderWriter_WriteHTML_NotFound(t *testing.T) {
	t.Parallel()

	rw := chtmltest.NewReaderWriter(t)
	resp := httptest.NewRecorder()

	rw.WriteHTML(resp, httptest.NewRequest(http.MethodGet, "/", nil), chtml.WriteHTMLParams{
		StatusCode: http.StatusNotFound,
		Data:       map[string]string{"user": "test"},
	})

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, "text/html", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), `not found`)
}
