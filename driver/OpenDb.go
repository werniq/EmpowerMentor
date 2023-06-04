package driver

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

// OpenDb opens a connection to the database, creates necessary tables if they don't exist
func OpenDb() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
