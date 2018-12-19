package csql

import (
	"context"
	"fmt"

	"github.com/tusharsoni/copper/clogger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgres dialect for gorm
	"go.uber.org/fx"
)

type gormDBParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    Config
	Logger    clogger.Logger
}

func newGormDB(p gormDBParams) (*gorm.DB, error) {
	conn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		p.Config.Host,
		p.Config.Port,
		p.Config.User,
		p.Config.Name,
	)
	p.Logger.Info("Connecting to database..", map[string]string{
		"connection": conn,
	})

	if p.Config.Password != "" {
		conn = fmt.Sprintf("%s password=%s", conn, p.Config.Password)
	}

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		p.Logger.Error("Failed to connect to database..", err)
		return nil, err
	}

	p.Lifecycle.Append(fx.Hook{
		OnStop: func(context context.Context) error {
			p.Logger.Info("Closing connection to database..", nil)
			return db.Close()
		},
	})

	db.LogMode(false)

	return db, nil
}
