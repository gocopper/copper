// +build csql_postgres

package csql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newDialect(dsn string) gorm.Dialector {
	return postgres.Open(dsn)
}
