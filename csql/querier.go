package csql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clogger"
	"github.com/jmoiron/sqlx"
)

// Querier provides a set of helpful methods to run database queries. It can be used to run parameterized queries
// and scan results into Go structs or slices.
type Querier interface {
	CtxWithTx(ctx context.Context) (context.Context, *sql.Tx, error)
	InTx(ctx context.Context, fn func(context.Context) error) error
	WithIn() Querier
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// NewQuerier returns a querier using the given database connection and the dialect
func NewQuerier(db *sql.DB, config Config, logger clogger.Logger) Querier {
	return &querier{
		db:      sqlx.NewDb(db, config.Dialect),
		dialect: config.Dialect,
		in:      false,
		logger:  logger,
	}
}

type querier struct {
	db      *sqlx.DB
	dialect string
	in      bool
	logger  clogger.Logger
}

func (q *querier) CtxWithTx(ctx context.Context) (context.Context, *sql.Tx, error) {
	return CtxWithTx(ctx, q.db.DB, q.dialect)
}

func (q *querier) InTx(ctx context.Context, fn func(context.Context) error) error {
	ctx, tx, err := CtxWithTx(ctx, q.db.DB, q.dialect)
	if err != nil {
		return cerrors.New(err, "failed to create context with database transaction", nil)
	}

	defer func() {
		// Try a rollback in a deferred function to account for panics
		err := tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			q.logger.Error("Failed to rollback database transaction", err)
			return
		}

		if err == nil {
			q.logger.Warn("Rolled back an unexpectedly open database transaction", nil)
		}
	}()

	err = fn(ctx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			q.logger.Error("Failed to rollback database transaction", err)
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return cerrors.New(err, "failed to commit database transaction", nil)
	}

	return nil
}

func (q *querier) WithIn() Querier {
	return &querier{
		db:      q.db,
		dialect: q.dialect,
		in:      true,
		logger:  q.logger,
	}
}

func (q *querier) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	query, args, err := q.mkQueryWithArgs(ctx, query, args)
	if err != nil {
		return err
	}

	return mustTxFromCtx(ctx).GetContext(ctx, dest, query, args...)
}

func (q *querier) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	query, args, err := q.mkQueryWithArgs(ctx, query, args)
	if err != nil {
		return err
	}

	return mustTxFromCtx(ctx).SelectContext(ctx, dest, query, args...)
}

func (q *querier) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	query, args, err := q.mkQueryWithArgs(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return mustTxFromCtx(ctx).ExecContext(ctx, query, args...)
}

func (q *querier) mkQueryWithArgs(ctx context.Context, query string, args []interface{}) (string, []interface{}, error) {
	var err error

	if q.in {
		query, args, err = sqlx.In(query, args...)
		if err != nil {
			return "", nil, cerrors.New(err, "failed to create IN query", nil)
		}
	}

	tx, err := txFromCtx(ctx)
	if err != nil {
		return "", nil, err
	}

	return tx.Rebind(query), args, nil
}
