package cauthtest

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gocopper/copper/chttp"

	"github.com/gocopper/copper/chttp/chttptest"

	"github.com/gocopper/copper/cauth"
	"github.com/gocopper/copper/clogger"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// NewHandler instantiates and returns a http.Handler with auth roter and middlewares suited for testing.
func NewHandler(t *testing.T) http.Handler {
	t.Helper()

	logger := clogger.New()
	rw := chttptest.NewReaderWriter(t)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{ // nolint: exhaustivestruct
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	assert.NoError(t, err)

	err = cauth.NewMigration(db).Run()
	assert.NoError(t, err)

	svc := cauth.NewSvc(cauth.NewRepo(db))

	setSessionMW := cauth.NewSetSessionMiddleware(svc, rw, logger)
	verifySessionMW := cauth.NewVerifySessionMiddleware(svc, rw, logger)

	router := cauth.NewRouter(cauth.NewRouterParams{
		Auth:      svc,
		RW:        rw,
		SessionMW: verifySessionMW,
		Logger:    logger,
	})

	handler := chttp.NewHandler(chttp.NewHandlerParams{
		Routers:           []chttp.Router{router},
		GlobalMiddlewares: []chttp.Middleware{setSessionMW},
	})

	return handler
}

// CreateNewUserSession creates a new user using the given router and returns the session created by it.
func CreateNewUserSession(t *testing.T, server *httptest.Server) *cauth.SessionResult {
	t.Helper()

	var session cauth.SessionResult

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

	err = json.Unmarshal(respBodyJ, &session)
	assert.NoError(t, err)

	return &session
}
