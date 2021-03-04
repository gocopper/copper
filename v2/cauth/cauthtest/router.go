package cauthtest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/cauth"
	"github.com/tusharsoni/copper/v2/chttp"
	"github.com/tusharsoni/copper/v2/clogger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// NewRouter instantiates and returns a router suited for testing.
func NewRouter(t *testing.T) *cauth.Router {
	t.Helper()

	logger := clogger.New()
	rw := chttp.NewJSONReaderWriter(logger)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{ // nolint: exhaustivestruct
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	assert.NoError(t, err)

	err = cauth.NewMigrator(db).Run()
	assert.NoError(t, err)

	svc := cauth.NewSvc(cauth.NewRepo(db))

	sessionMW := cauth.NewVerifySessionMiddleware(svc, rw, logger)

	return cauth.NewRouter(cauth.NewRouterParams{
		Auth:      svc,
		RW:        rw,
		SessionMW: sessionMW,
		Logger:    logger,
	})
}

// CreateNewUserSession creates a new user using the given router and returns the session created by it.
func CreateNewUserSession(t *testing.T, router *cauth.Router) *cauth.SessionResult {
	t.Helper()

	var session cauth.SessionResult

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

	err = json.Unmarshal(respBodyJ, &session)
	assert.NoError(t, err)

	return &session
}
