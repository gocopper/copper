package csql_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gocopper/copper/clifecycle/clifecycletest"
	"github.com/gocopper/copper/clogger"
	"github.com/gocopper/copper/csql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestQuerier_Read(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	_, err = db.Exec("create table people (name text);insert into people (name) values ('test');")
	assert.NoError(t, err)

	querier := csql.NewQuerier(db, clifecycletest.New(), csql.Config{Dialect: "sqlite3"}, clogger.NewNoop())

	ctx, _, err := csql.CtxWithTx(context.Background(), db, "sqlite3")
	assert.NoError(t, err)

	t.Run("get", func(t *testing.T) {
		t.Parallel()

		var dest struct {
			Name string
		}

		err = querier.Get(ctx, &dest, "select * from people")
		assert.NoError(t, err)

		assert.Equal(t, "test", dest.Name)
	})

	t.Run("select", func(t *testing.T) {
		t.Parallel()

		var dest []struct {
			Name string
		}

		err = querier.Select(ctx, &dest, "select * from people")
		assert.NoError(t, err)

		assert.Equal(t, 1, len(dest))
		assert.Equal(t, "test", dest[0].Name)
	})

	t.Run("select in", func(t *testing.T) {
		t.Parallel()

		var dest []struct {
			Name string
		}

		err = querier.WithIn().Select(ctx, &dest, "select * from people where name in (?)", []string{"test"})
		assert.NoError(t, err)

		assert.Equal(t, 1, len(dest))
		assert.Equal(t, "test", dest[0].Name)
	})
}

func TestQuerier_Exec(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	_, err = db.Exec("create table people (name text);insert into people (name) values ('test');")
	assert.NoError(t, err)

	querier := csql.NewQuerier(db, clifecycletest.New(), csql.Config{Dialect: "sqlite3"}, clogger.NewNoop())

	ctx, _, err := csql.CtxWithTx(context.Background(), db, "sqlite3")
	assert.NoError(t, err)

	res, err := querier.Exec(ctx, "delete from people")
	assert.NoError(t, err)

	n, err := res.RowsAffected()
	assert.NoError(t, err)

	assert.Equal(t, int64(1), n)
}

func TestQuerier_DirectQueries(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	_, err = db.Exec("create table people (name text);")
	assert.NoError(t, err)

	querier := csql.NewQuerier(db, clifecycletest.New(), csql.Config{Dialect: "sqlite3"}, clogger.NewNoop())
	ctx := context.Background()

	// Queries run directly on DB without transaction
	_, err = querier.Exec(ctx, "insert into people (name) values (?)", "alice")
	assert.NoError(t, err)

	_, err = querier.Exec(ctx, "insert into people (name) values (?)", "bob")
	assert.NoError(t, err)

	var dest struct {
		Name string
	}
	err = querier.Get(ctx, &dest, "select * from people where name = ?", "alice")
	assert.NoError(t, err)
	assert.Equal(t, "alice", dest.Name)

	// Verify data was committed
	var count int
	err = db.QueryRow("select count(*) from people").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestQuerier_CtxWithoutTx(t *testing.T) {
	t.Parallel()

	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	// Verify CtxWithoutTx removes transaction from context
	ctx, _, err := csql.CtxWithTx(context.Background(), db, "sqlite3")
	assert.NoError(t, err)

	// Original context has transaction
	_, err = csql.TxFromCtx(ctx)
	assert.NoError(t, err, "original context should have transaction")

	// CtxWithoutTx removes transaction
	ctxNoTx := csql.CtxWithoutTx(ctx)
	_, err = csql.TxFromCtx(ctxNoTx)
	assert.Error(t, err, "CtxWithoutTx should remove transaction from context")
	assert.Contains(t, err.Error(), "no database transaction in the context")
}
