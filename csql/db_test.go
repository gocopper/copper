package csql_test

import (
	"testing"

	"github.com/gocopper/copper"
	"github.com/gocopper/copper/cconfig"
	"github.com/gocopper/copper/cconfig/cconfigtest"
	"github.com/gocopper/copper/clogger"
	"github.com/gocopper/copper/csql"
	"github.com/stretchr/testify/assert"
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
