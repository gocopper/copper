package csql

import (
	"context"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/clogger"
)

// DBTxnMiddleware provides a middleware that wraps the http request in a database transaction. If the response status
// code is 2xx, the transaction is committed, else a rollback is performed.
type DBTxnMiddleware interface {
	WrapInTxn(next http.Handler) http.Handler
}

type dbTxnMiddleware struct {
	db     *gorm.DB
	logger clogger.Logger
}

func newDBTxnMiddleware(db *gorm.DB, logger clogger.Logger) DBTxnMiddleware {
	return &dbTxnMiddleware{
		db:     db,
		logger: logger,
	}
}

func (m *dbTxnMiddleware) WrapInTxn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := m.db.Begin()
		ctx := context.WithValue(r.Context(), connCtxKey, txn)

		rw := &txnrw{
			internal: w,
			db:       txn,
			logger:   m.logger,
		}

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

type txnrw struct {
	internal http.ResponseWriter
	db       *gorm.DB
	logger   clogger.Logger
}

func (w *txnrw) Header() http.Header {
	return w.internal.Header()
}

func (w *txnrw) Write(b []byte) (int, error) {
	return w.internal.Write(b)
}

func (w *txnrw) WriteHeader(statusCode int) {
	if statusCode < 200 || statusCode > 299 {
		w.db.Rollback()
		w.internal.WriteHeader(statusCode)
		return
	}

	err := w.db.Commit().Error
	if err != nil {
		w.logger.Error("Failed to commit db transaction", err)
		w.internal.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.internal.WriteHeader(statusCode)
}
