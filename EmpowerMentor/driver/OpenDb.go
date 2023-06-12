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

	stmt := `CREATE TABLE if not exists requested_challenges
				(
				id BIGSERIAL PRIMARY KEY, 
				user_id INTEGER NOT NULL, 
				challenge_id INTEGER NOT NULL, 
				challenge TEXT NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);`

	_, err = db.Exec(stmt)
	if err != nil {
		return nil, err
	}

	stmt = `CREATE TABLE IF NOT EXISTS week (
			  id SERIAL PRIMARY KEY,
			  user_id BIGINT NOT NULL,
			  monday JSONB,
			  tuesday JSONB,
			  wednesday JSONB,
			  thursday JSONB,
			  friday JSONB,
			  saturday JSONB,
			  sunday JSONB
			);
				`

	_, err = db.Exec(stmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}
