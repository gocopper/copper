package csql_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gocopper/copper/clifecycle/clifecycletest"
	"github.com/gocopper/copper/clogger"
	"github.com/gocopper/copper/csql"
	"github.com/stretchr/testify/assert"
)

func TestCtxWithTx(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	querier := csql.NewQuerier(db, clifecycletest.New(), csql.Config{Dialect: "sqlite3"}, clogger.NewNoop())

	ctx, err := querier.CtxWithTx(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, ctx)
}
