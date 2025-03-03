package chttp_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocopper/copper/chttp"
	"github.com/gocopper/copper/clogger"
	"github.com/stretchr/testify/assert"
)

func TestJSONReaderWriter_ReadJSON(t *testing.T) {
	t.Parallel()

	var body struct {
		Key string `json:"key"`
	}

	rw := chttp.NewJSONReaderWriter(chttp.Config{}, clogger.NewNoop())

	ok := rw.ReadJSON(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(`{"key": "value"}`))),
		&body,
	)

	assert.True(t, ok)
	assert.Equal(t, "value", body.Key)
}

func TestJSONReaderWriter_ReadJSON_Invalid_Body(t *testing.T) {
	t.Parallel()

	var body struct {
		Key string `json:"key"`
	}

	rw := chttp.NewJSONReaderWriter(chttp.Config{}, clogger.NewNoop())
	resp := httptest.NewRecorder()

	ok := rw.ReadJSON(
		resp,
		httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(`{ invalid json }`))),
		&body,
	)

	assert.False(t, ok)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), "invalid body json")
}

func TestJSONReaderWriter_ReadJSON_Empty_Body(t *testing.T) {
	t.Parallel()

	var body struct {
		Key string `json:"key"`
	}

	rw := chttp.NewJSONReaderWriter(chttp.Config{}, clogger.NewNoop())
	resp := httptest.NewRecorder()

	ok := rw.ReadJSON(
		resp,
		httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(``))),
		&body,
	)

	assert.False(t, ok)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), "empty body")
}

func TestJSONReaderWriter_ReadJSON_Validator(t *testing.T) {
	t.Parallel()

	var body struct {
		Key string `json:"key" valid:"email"`
	}

	rw := chttp.NewJSONReaderWriter(chttp.Config{}, clogger.NewNoop())
	resp := httptest.NewRecorder()

	ok := rw.ReadJSON(
		resp,
		httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(`{"key": "value"}`))),
		&body,
	)

	assert.False(t, ok)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestJSONReaderWriter_WriteJSON_Data(t *testing.T) {
	t.Parallel()

	rw := chttp.NewJSONReaderWriter(chttp.Config{}, clogger.NewNoop())
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

func TestJSONReaderWriter_WriteJSON_StatusCode(t *testing.T) {
	t.Parallel()

	rw := chttp.NewJSONReaderWriter(chttp.Config{}, clogger.NewNoop())
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

func TestJSONReaderWriter_WriteJSON_Error(t *testing.T) {
	t.Parallel()

	rw := chttp.NewJSONReaderWriter(chttp.Config{}, clogger.NewNoop())
	resp := httptest.NewRecorder()

	rw.WriteJSON(resp, chttp.WriteJSONParams{
		StatusCode: http.StatusBadRequest,
		Data:       errors.New("test-err"), //nolint:goerr113
	})

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), `{"error":"test-err"}`)
}

func TestJSONReaderWriter_WriteJSON_NilData(t *testing.T) {
	t.Parallel()

	rw := chttp.NewJSONReaderWriter(chttp.Config{}, clogger.NewNoop())
	resp := httptest.NewRecorder()

	rw.WriteJSON(resp, chttp.WriteJSONParams{
		StatusCode: http.StatusOK,
		Data:       nil,
	})

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("content-type"))
	assert.Empty(t, resp.Body.String())
}

func TestJSONReaderWriter_Unauthorized(t *testing.T) {
	t.Parallel()

	rw := chttp.NewJSONReaderWriter(chttp.Config{}, clogger.NewNoop())
	resp := httptest.NewRecorder()

	rw.Unauthorized(resp)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("content-type"))
	assert.Contains(t, resp.Body.String(), `{"error":"unauthorized"}`)
}
