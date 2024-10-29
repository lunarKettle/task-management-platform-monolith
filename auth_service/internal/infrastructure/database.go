package infrastructure

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	connection *sql.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) OpenConnetion(connectionString string) (err error) {
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
