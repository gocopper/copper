package csql

import (
	"context"
	"fmt"

	"github.com/tusharsoni/copper/clogger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Postgres dialect for gorm
	"go.uber.org/fx"
)

type GormDBParams struct {
	fx.In

	Config    Config
	Logger    clogger.Logger
	Lifecycle fx.Lifecycle `optional:"true"`
}

func NewGormDB(p GormDBParams) (*gorm.DB, error) {
	conn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=disable",
		p.Config.Host,
		p.Config.Port,
		p.Config.User,
		p.Config.Name,
	)

	p.Logger.WithTags(map[string]interface{}{
		"connection": conn,
	}).Info("Connecting to database...")

	if p.Config.Password != "" {
		conn = fmt.Sprintf("%s password=%s", conn, p.Config.Password)
	}

	db, err := gorm.Open("postgres", conn)
	if err != nil {
		p.Logger.Error("Failed to connect to database..", err)
		return nil, err
	}

	if p.Lifecycle != nil {
		p.Lifecycle.Append(fx.Hook{
			OnStop: func(context context.Context) error {
				p.Logger.Info("Closing connection to database..")
				return db.Close()
			},
		})
	}

	db.LogMode(false)

	return db, nil
}
