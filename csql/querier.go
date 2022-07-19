package csql

import (
	"context"
	"database/sql"

	"github.com/gocopper/copper/cerrors"
	"github.com/jmoiron/sqlx"
)

// Querier provides a set of helpful methods to run database queries. It can be used to run parameterized queries
// and scan results into Go structs or slices.
type Querier interface {
	WithIn() Querier
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// NewQuerier returns a querier using the given database connection and the dialect
func NewQuerier(db *sql.DB, config Config) Querier {
	return &querier{
		db: sqlx.NewDb(db, config.Dialect),
		in: false,
	}
}

type querier struct {
	db *sqlx.DB
	in bool
}

func (q *querier) WithIn() Querier {
	return &querier{
		db: q.db,
		in: true,
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
