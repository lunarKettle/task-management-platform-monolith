package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	connection *sql.DB
}

func (db *Database) OpenConnetion() (err error) {
	connectionString := "user=admin password=admin dbname=project_management_service_db sslmode=disable"
	db.connection, err = sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("failed to open connection: %w", err)
	}

	if err := db.CheckConnection(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}

func (db *Database) CloseConnection() error {
	return db.connection.Close()
}

func (db *Database) CheckConnection() error {
	if err := db.connection.Ping(); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}
