// +build csql_sqlite

package csql

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newDialect(dsn string) gorm.Dialector {
	return sqlite.Open(dsn)
}
