package postgresql

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func Migrate(fs *embed.FS, cfg *Config) error {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	source, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations source %v", err)
	}

	instance, err := migrate.NewWithSourceInstance("iofs", source, dbUrl)
	if err != nil {
		return fmt.Errorf("failed to initialization the migrations instance %v", err)
	}

	err = instance.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
