package csql

import (
	"context"

	"gorm.io/gorm"
)

type ctxKey string

const connCtxKey = ctxKey("csql/tx")

// GetConn returns a db connection from the context or the given default connection if context is empty.
func GetConn(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx, ok := ctx.Value(connCtxKey).(*gorm.DB)
	if !ok {
		return db.WithContext(ctx)
	}

	return tx.WithContext(ctx)
}
