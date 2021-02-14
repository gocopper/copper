package chttp_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/chttp"
	"github.com/tusharsoni/copper/v2/clogger"
)

func TestNewJSONReaderWriter(t *testing.T) {
	t.Parallel()

	rw := chttp.NewJSONReaderWriter(clogger.NewNoop())

	_, ok := rw.(chttp.ReaderWriter)

	assert.NotNil(t, rw)
	assert.True(t, ok)
}

func TestJSONRW_Read(t *testing.T) {
	t.Parallel()

	var body struct {
		Key string `json:"key"`
	}

	rw := chttp.NewJSONReaderWriter(clogger.NewNoop())

	ok := rw.Read(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(`{"key": "value"}`))),
		&body,
	)

	assert.True(t, ok)
	assert.Equal(t, "value", body.Key)
}

func TestJSONRW_Read_Invalid_Body(t *testing.T) {
	t.Parallel()

	var body struct {
		Key string `json:"key"`
	}

	rw := chttp.NewJSONReaderWriter(clogger.NewNoop())

	ok := rw.Read(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(`{ invalid json }`))),
		&body,
	)

	assert.False(t, ok)
}

func TestJSONRW_Read_Validator(t *testing.T) {
	t.Parallel()

	var body struct {
		Key string `json:"key" valid:"email"`
	}

	rw := chttp.NewJSONReaderWriter(clogger.NewNoop())

	ok := rw.Read(
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(`{"key": "value"}`))),
		&body,
	)

	assert.False(t, ok)
}

func TestJSONRW_OK(t *testing.T) {
	t.Parallel()

	rw := chttp.NewJSONReaderWriter(clogger.NewNoop())
	resp := httptest.NewRecorder()

	rw.OK(resp, map[string]string{
		"key": "val",
	})

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("content-type"))
	assert.Equal(t, `{"key":"val"}`, resp.Body.String())
}

func TestJSONRW_Created(t *testing.T) {
	t.Parallel()

	rw := chttp.NewJSONReaderWriter(clogger.NewNoop())
	resp := httptest.NewRecorder()

	rw.Created(resp, map[string]string{
		"key": "val",
	})

	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("content-type"))
	assert.Equal(t, `{"key":"val"}`, resp.Body.String())
}

func TestJSONRW_InternalErr(t *testing.T) {
	t.Parallel()

	rw := chttp.NewJSONReaderWriter(clogger.NewNoop())
	resp := httptest.NewRecorder()

	rw.InternalErr(resp)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestJSONRW_Unauthorized(t *testing.T) {
	t.Parallel()

	rw := chttp.NewJSONReaderWriter(clogger.NewNoop())
	resp := httptest.NewRecorder()

	rw.Unauthorized(resp)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}

func TestJSONRW_Forbidden(t *testing.T) {
	t.Parallel()

	rw := chttp.NewJSONReaderWriter(clogger.NewNoop())
	resp := httptest.NewRecorder()

	rw.Forbidden(resp)

	assert.Equal(t, http.StatusForbidden, resp.Code)
}

func TestJSONRW_BadRequest(t *testing.T) {
	t.Parallel()

	rw := chttp.NewJSONReaderWriter(clogger.NewNoop())
	resp := httptest.NewRecorder()

	rw.BadRequest(resp, errors.New("test-err")) //nolint:goerr113

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, "application/json", resp.Header().Get("content-type"))
	assert.Equal(t, `{"error":"test-err"}`, resp.Body.String())
}
