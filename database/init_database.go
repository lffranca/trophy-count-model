package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// InitDatabase InitDatabase
func InitDatabase() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=admin dbname=trophy password=admin sslmode=disable")
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)

	return db, nil
}
