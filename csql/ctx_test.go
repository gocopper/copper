package csql_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gocopper/copper/csql"
	"github.com/stretchr/testify/assert"
)

func TestCtxWithTx(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	ctx, tx1, err := csql.CtxWithTx(context.Background(), db, "sqlite3")
	assert.NoError(t, err)

	tx2, err := csql.TxFromCtx(ctx)
	assert.NoError(t, err)

	assert.Equal(t, tx1, tx2)
}
