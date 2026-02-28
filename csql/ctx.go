package csql

import (
	"context"
	"database/sql"

	"github.com/gocopper/copper/cerrors"
	"github.com/jmoiron/sqlx"
)

type ctxKey string

const connCtxKey = ctxKey("csql/*sqlx.Tx")

func ctxWithTx(parentCtx context.Context, db *sql.DB, dialect string) (context.Context, error) {
	tx, err := sqlx.NewDb(db, dialect).Beginx()
	if err != nil {
		return nil, cerrors.New(err, "failed to begin db transaction", map[string]any{
			"dialect": dialect,
		})
	}

	return context.WithValue(parentCtx, connCtxKey, tx), nil
}

func txFromCtx(ctx context.Context) (*sqlx.Tx, error) {
	tx, ok := ctx.Value(connCtxKey).(*sqlx.Tx)
	if !ok {
		return nil, cerrors.New(nil, "no database transaction in the context", nil)
	}

	return tx, nil
}

func txFromCtxOrNil(ctx context.Context) *sqlx.Tx {
	tx, _ := ctx.Value(connCtxKey).(*sqlx.Tx)
	return tx
}

// contextWithoutTx wraps a context and shadows the transaction key while preserving all other values.
type contextWithoutTx struct {
	context.Context
}

func (c *contextWithoutTx) Value(key any) any {
	if key == connCtxKey {
		return nil // Remove transaction
	}
	return c.Context.Value(key) // Preserve all other values
}

// CtxWithoutTx returns a context without the transaction, forcing queries to use auto-commit behavior.
// Preserves all other context values (tracing, auth, deadlines, etc.).
// Useful when you need to persist data regardless of whether the parent transaction rolls back
// (e.g., audit logs, metrics).
func CtxWithoutTx(ctx context.Context) context.Context {
	return &contextWithoutTx{Context: ctx}
}
