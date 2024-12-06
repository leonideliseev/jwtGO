package app

import (
	"os"

	"github.com/leonideliseev/jwtGO/pkg/postgresql"
	"github.com/leonideliseev/jwtGO/schema"
	"github.com/spf13/viper"
)

func (a *App) initDBConn() error {
	config := postgresql.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	configTest := config
	configTest.DBName = "postgres"
	connTest, err := postgresql.ConnWithPgxPool(configTest)
	if err != nil {
		return err
	}
	defer connTest.Close()

	err = postgresql.CreateDatabaseIfNotExists(connTest, viper.GetString("db.dbname"))
	if err != nil {
		panic(err)
	}

	conn, err := postgresql.ConnWithPgxPool(config)
	if err != nil {
		return err
	}

	err = postgresql.Migrate(&schema.DB, &config)
	if err != nil {
		panic(err)
	}

	a.conn = conn
	return nil
}
