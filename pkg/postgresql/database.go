package postgresql

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDatabaseIfNotExists(conn *pgxpool.Pool, dbName string) error {
	var exists bool
	err := conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists %v", err)
	}

	if !exists {
		_, err := conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database %v", err)
		}
	}

	log.Printf("database created %s", dbName)

	return nil
}
