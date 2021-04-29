package chttp_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocopper/copper/chttp/chttptest"

	"github.com/gocopper/copper/chttp"
	"github.com/stretchr/testify/assert"
)

func TestReaderWriter_ReadJSON(t *testing.T) {
	t.Parallel()

	var body struct {
		Key string `json:"key"`
	}

	rw := chttptest.NewReaderWriter(t)

	ok := rw.ReadJSON(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(`{"key": "value"}`))),
		&body,
	)

	assert.True(t, ok)
	assert.Equal(t, "value", body.Key)
}

func TestReaderWriter_ReadJSON_Invalid_Body(t *testing.T) {
	t.Parallel()

	var body struct {
		Key string `json:"key"`
	}

	rw := chttptest.NewReaderWriter(t)

	ok := rw.ReadJSON(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(`{ invalid json }`))),
		&body,
	)

	assert.False(t, ok)
}

func TestReaderWriter_ReadJSON_Validator(t *testing.T) {
	t.Parallel()

	var body struct {
		Key string `json:"key" valid:"email"`
	}

	rw := chttptest.NewReaderWriter(t)

	ok := rw.ReadJSON(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(`{"key": "value"}`))),
		&body,
	)

	assert.False(t, ok)
}

func TestReaderWriter_WriteJSON_Data(t *testing.T) {
	t.Parallel()

	rw := chttptest.NewReaderWriter(t)
	resp := httptest.NewRecorder()

	rw.WriteJSON(resp, chttp.WriteJSONParams{
		Data: map[string]string{
			"key": "val",
		},
	})

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), `{"key":"val"}`)
}

func TestReaderWriter_WriteJSON_StatusCode(t *testing.T) {
	t.Parallel()

	rw := chttptest.NewReaderWriter(t)
	resp := httptest.NewRecorder()

	rw.WriteJSON(resp, chttp.WriteJSONParams{
		StatusCode: http.StatusCreated,
		Data: map[string]string{
			"key": "val",
		},
	})

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), `{"key":"val"}`)
}

func TestReaderWriter_WriteJSON_Error(t *testing.T) {
	t.Parallel()

	rw := chttptest.NewReaderWriter(t)
	resp := httptest.NewRecorder()

	rw.WriteJSON(resp, chttp.WriteJSONParams{
		StatusCode: http.StatusBadRequest,
		Data:       errors.New("test-err"), //nolint:goerr113
	})

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), `{"error":"test-err"}`)
}

func TestReaderWriter_WriteHTML(t *testing.T) {
	t.Parallel()

	rw := chttptest.NewReaderWriter(t)
	resp := httptest.NewRecorder()

	rw.WriteHTML(resp, httptest.NewRequest(http.MethodGet, "/", nil), chttp.WriteHTMLParams{
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

	rw := chttptest.NewReaderWriter(t)
	resp := httptest.NewRecorder()

	rw.WriteHTML(resp, httptest.NewRequest(http.MethodGet, "/", nil), chttp.WriteHTMLParams{
		StatusCode: http.StatusNotFound,
		Data:       map[string]string{"user": "test"},
	})

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Equal(t, "text/html", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), `not found`)
}
