package csql

import (
	"context"
	"database/sql"

	"github.com/gocopper/copper/cerrors"
	"github.com/jmoiron/sqlx"
)

type ctxKey string

const connCtxKey = ctxKey("csql/*sqlx.Tx")

// CtxWithTx creates a context with a new database transaction. Any queries run using Querier will be run within
// this transaction.
func CtxWithTx(parentCtx context.Context, db *sql.DB, dialect string) (context.Context, *sql.Tx, error) {
	tx, err := sqlx.NewDb(db, dialect).Beginx()
	if err != nil {
		return nil, nil, cerrors.New(err, "failed to begin db transaction", map[string]interface{}{
			"dialect": dialect,
		})
	}

	return context.WithValue(parentCtx, connCtxKey, tx), tx.Tx, nil
}

// TxFromCtx returns an existing transaction from the context. This method should be called with context created
// using CtxWithTx.
func TxFromCtx(ctx context.Context) (*sql.Tx, error) {
	tx, err := txFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return tx.Tx, nil
}

func txFromCtx(ctx context.Context) (*sqlx.Tx, error) {
	tx, ok := ctx.Value(connCtxKey).(*sqlx.Tx)
	if !ok {
		return nil, cerrors.New(nil, "no database transaction in the context", nil)
	}

	return tx, nil
}

func mustTxFromCtx(ctx context.Context) *sqlx.Tx {
	tx, err := txFromCtx(ctx)
	if err != nil {
		panic(err)
	}

	return tx
}
