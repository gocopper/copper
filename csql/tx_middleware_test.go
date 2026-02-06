package csql_test

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gocopper/copper/clifecycle/clifecycletest"
	"github.com/gocopper/copper/clogger"
	"github.com/gocopper/copper/csql"
	"github.com/stretchr/testify/assert"
)

func TestTxMiddleware_Handle_Commit(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	_, err = db.Exec("create table people (name text)")
	assert.NoError(t, err)

	var (
		logger  = clogger.NewNoop()
		config  = csql.Config{Dialect: "sqlite3"}
		lc      = clifecycletest.New()
		querier = csql.NewQuerier(db, lc, config, logger)
		mw      = csql.NewTxMiddleware(db, querier, config, logger)
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", nil)
	assert.NoError(t, err)

	mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := querier.Exec(r.Context(), "insert into people (name) values ('test')")
		assert.NoError(t, err)
	})).ServeHTTP(httptest.NewRecorder(), req)

	rows, err := db.Query("select * from people")
	assert.NoError(t, err)
	assert.NoError(t, rows.Err())

	assert.True(t, rows.Next())
}

func TestTxMiddleware_Handle_Rollback(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	_, err = db.Exec("create table people (name text)")
	assert.NoError(t, err)

	var (
		logger  = clogger.NewNoop()
		config  = csql.Config{Dialect: "sqlite3"}
		lc      = clifecycletest.New()
		querier = csql.NewQuerier(db, lc, config, logger)
		mw      = csql.NewTxMiddleware(db, querier, config, logger)
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "/", nil)
	assert.NoError(t, err)

	mw.Handle(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := querier.Exec(r.Context(), "insert into people (name) values ('test')")
		assert.NoError(t, err)

		w.WriteHeader(http.StatusInternalServerError)
	})).ServeHTTP(httptest.NewRecorder(), req)

	rows, err := db.Query("select * from people")
	assert.NoError(t, err)
	assert.NoError(t, rows.Err())

	assert.False(t, rows.Next())
}
