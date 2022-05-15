package csql_test

import (
	"testing"

	"github.com/gocopper/copper/clifecycle"
	"github.com/gocopper/copper/clogger"
	"github.com/gocopper/copper/csql"
	"github.com/stretchr/testify/assert"
)

func TestNewDBConnection(t *testing.T) {
	t.Parallel()

	var (
		logger = clogger.New()
		lc     = clifecycle.New()
	)

	db, err := csql.NewDBConnection(lc, csql.Config{
		Dialect: "sqlite",
		DSN:     ":memory:",
	}, logger)
	assert.NoError(t, err)

	sqlDB, err := db.DB()
	assert.NoError(t, err)

	assert.NoError(t, sqlDB.Ping())

	lc.Stop(logger)

	assert.Error(t, sqlDB.Ping())
}
