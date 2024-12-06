package app

import (
	"github.com/leonideliseev/jwtGO/pkg/postgresql"
	"github.com/leonideliseev/jwtGO/schema"
)

func (a *App) initDBConn(cfg postgresql.Config) error {
	configTest := cfg
	configTest.DBName = "postgres"
	connTest, err := postgresql.ConnWithPgxPool(configTest)
	if err != nil {
		return err
	}
	defer connTest.Close()

	err = postgresql.CreateDatabaseIfNotExists(connTest, cfg.DBName)
	if err != nil {
		return err
	}

	conn, err := postgresql.ConnWithPgxPool(cfg)
	if err != nil {
		return err
	}

	err = postgresql.Migrate(&schema.DB, &cfg)
	if err != nil {
		return err
	}

	a.conn = conn
	return nil
}
