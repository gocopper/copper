package csql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper"
	"github.com/tusharsoni/copper/cconfig"
	"github.com/tusharsoni/copper/cconfig/cconfigtest"
	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/csql"
)

func TestNewDBConnection(t *testing.T) {
	t.Parallel()

	logger := clogger.New()
	lc := copper.NewLifecycle(logger)
	config, err := cconfig.New(cconfigtest.SetupDirWithConfigs(t, `
[csql]
dsn = ":memory:"
`, ""), ".", "test")
	assert.NoError(t, err)

	db, err := csql.NewDBConnection(lc, config, logger)
	assert.NoError(t, err)

	sqlDB, err := db.DB()
	assert.NoError(t, err)

	assert.NoError(t, sqlDB.Ping())

	lc.Stop()

	assert.Error(t, sqlDB.Ping())
}
