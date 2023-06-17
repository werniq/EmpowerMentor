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

	stmt := `
		CREATE TABLE IF NOT EXISTS nutrients (
			id bigserial not null primary key,
			calories double precision,
			protein double precision,
			fat double precision,
			carbohydrates double precision
		);


		CREATE TABLE IF NOT EXISTS meals (
			id bigserial not null primary key,
			image_type text,
			title varchar(255),
			ready_in_minutes int,
			servings int,
			source_url text
		);

		CREATE TABLE IF NOT EXISTS monday (
			id bigserial not null primary key,
			user_id bigint not null,
			meal_ids bigint[] not null,
			nutrient_id bigint not null,
			foreign key (nutrient_id) references nutrients(id)
		);

		CREATE TABLE IF NOT EXISTS tuesday (
			id bigserial not null primary key,
			user_id bigint not null,
			meal_ids bigint[] not null,
			nutrient_id bigint not null,
			foreign key (nutrient_id) references nutrients(id)
		);

		CREATE TABLE IF NOT EXISTS wednesday (
			id bigserial not null primary key,
			user_id bigint not null,
			meal_ids bigint[] not null,
			nutrient_id bigint not null,
			foreign key (nutrient_id) references nutrients(id)
		);

		CREATE TABLE IF NOT EXISTS thursday (
			id bigserial not null primary key,
			user_id bigint not null,
			meal_ids bigint[] not null,
			nutrient_id bigint not null,
			foreign key (nutrient_id) references nutrients(id)
		);

		CREATE TABLE IF NOT EXISTS friday (
			id bigserial not null primary key,
			user_id bigint not null,
			meal_ids bigint[] not null,
			nutrient_id bigint not null,
			foreign key (nutrient_id) references nutrients(id)
		);

		CREATE TABLE IF NOT EXISTS saturday (
			id bigserial not null primary key,
			user_id bigint not null,
			meal_ids bigint[] not null,
			nutrient_id bigint not null,
			foreign key (nutrient_id) references nutrients(id)
		);

		CREATE TABLE IF NOT EXISTS sunday (
			id bigserial not null primary key,
			user_id bigint not null,
			meal_ids bigint[] not null,
			nutrient_id bigint not null,
			foreign key (nutrient_id) references nutrients(id)
		);

		CREATE TABLE IF NOT EXISTS week (
			id bigserial not null primary key,
			user_id bigint not null,
			monday bigint not null,
			tuesday bigint not null,
			wednesday bigint not null,
			thursday bigint not null,
			friday bigint not null,
			saturday bigint not null,
			sunday bigint not null,
			foreign key (monday) references monday(id),
			foreign key (tuesday) references tuesday(id),
			foreign key (wednesday) references wednesday(id),
			foreign key (thursday) references thursday(id),
			foreign key (friday) references friday(id),
			foreign key (saturday) references saturday(id),
			foreign key (sunday) references sunday(id)
		);

		CREATE TABLE IF NOT EXISTS requested_challenges (
			id BIGSERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			challenge_id INTEGER NOT NULL,
			challenge TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err = db.Exec(stmt)
	if err != nil {
		return nil, err
	}

	return db, nil
}
