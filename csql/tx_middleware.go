package csql

import (
	"bufio"
	"database/sql"
	"errors"
	"github.com/gocopper/copper/cerrors"
	"github.com/gocopper/copper/clogger"
	"net"
	"net/http"
)

// NewTxMiddleware creates a new TxMiddleware
func NewTxMiddleware(db *sql.DB, querier Querier, config Config, logger clogger.Logger) *TxMiddleware {
	return &TxMiddleware{
		db:      db,
		querier: querier,
		config:  config,
		logger:  logger,
	}
}

// TxMiddleware is a chttp.Middleware that wraps an HTTP request in a database transaction. If the request succeeds
// (i.e. 2xx or 3xx response code), the transaction is committed. Else, the transaction is rolled back.
type TxMiddleware struct {
	db      *sql.DB
	querier Querier
	config  Config
	logger  clogger.Logger
}

// Handle implements the chttp.Middleware interface. See TxMiddleware
func (m *TxMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, tx, err := CtxWithTx(r.Context(), m.db, m.config.Dialect)
		if err != nil {
			m.logger.Error("Failed to create context with database transaction", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		defer func() {
			// Try a rollback in a deferred function to account for panics
			err := tx.Rollback()
			if err != nil && !errors.Is(err, sql.ErrTxDone) {
				m.logger.Error("Failed to rollback database transaction", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if err == nil {
				m.logger.Warn("Rolled back an unexpectedly open database transaction", nil)
			}
		}()

		next.ServeHTTP(&txnrw{
			internal: w,
			tx:       tx,
			querier:  m.querier,
			logger:   m.logger,
		}, r.WithContext(ctx))

		err = m.querier.CommitTx(tx)
		if err != nil {
			m.logger.Error("Failed to commit database transaction", err)
			return
		}
	})
}

type txnrw struct {
	internal http.ResponseWriter
	tx       *sql.Tx
	querier  Querier
	logger   clogger.Logger
}

func (w *txnrw) Header() http.Header {
	return w.internal.Header()
}

func (w *txnrw) Write(b []byte) (int, error) {
	err := w.querier.CommitTx(w.tx)
	if err != nil {
		return 0, cerrors.New(err, "failed to commit database transaction", nil)
	}

	return w.internal.Write(b)
}

func (w *txnrw) WriteHeader(statusCode int) {
	const MinErrStatusCode = 400

	if statusCode >= MinErrStatusCode {
		err := w.tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			w.logger.WithTags(map[string]interface{}{
				"originalStatusCode": statusCode,
			}).Error("Failed to rollback database transaction", err)
			w.internal.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.internal.WriteHeader(statusCode)
		return
	}

	err := w.querier.CommitTx(w.tx)
	if err != nil {
		w.internal.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.internal.WriteHeader(statusCode)
}

func (w *txnrw) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.internal.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("internal response writer is not http.Hijacker")
	}

	return h.Hijack()
}
