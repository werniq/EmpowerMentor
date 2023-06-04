package models

import "time"

type User struct {
	Id                        string    `json:"id"`
	UserId                    string    `json:"user_id"`
	Sex                       bool      `json:"sex"`
	Age                       int       `json:"age"`
	Weight                    float32   `json:"weight"`
	Height                    float32   `json:"height"`
	PreferredPhysicalActivity string    `json:"preferred_physical_activity"`
	Username                  string    `json:"username"`
	WorkoutCount              int       `json:"workout_count"`
	BooksCount                int       `json:"books_count"`
	PreferringSupplements     string    `json:"preferring_supplements"`
	HabitsAcquired            int       `json:"habits_acquired"`
	Subscribed                bool      `json:"subscribed"`
	SubscribedAt              time.Time `json:"subscribed_at"`
	SubscribedForXDays        int       `json:"subscribed_for_x_days"`
}

// UserExists checks if user exists in the database
func (m *DatabaseModel) UserExists(userId int64) (bool, error) {
	stmt := `
			SELECT 
				COUNT(*)
			FROM
			    user_bot_configuration
			WHERE
			    user_id = $1;`

	var count int
	err := m.DB.QueryRow(stmt, userId).Scan(&count)
	if err != nil {
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// consider add customer function instead of add user

// IsAdmin function verifies if user is admin
func (m *DatabaseModel) IsAdmin(userId int64) error {
	stmt := `
		SELECT
			id
		FROM 
		    admins 
		WHERE 
		    user_id = $1
			`

	var id int
	err := m.DB.QueryRow(stmt).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (m *DatabaseModel) AddUser(userId int64) error {
	return nil
}

func (m *DatabaseModel) DeleteUser(userId int64) error {
	return nil
}
func (m *DatabaseModel) UpdateUser(userId int64) error {
	return nil
}

// AddAdmin function adds admin to database
func (m *DatabaseModel) AddAdmin(userID int64, username string) error {
	stmt := "INSERT INTO admins(user_id, username) VALUES ($1, $2)"

	_, err := m.DB.Exec(stmt, userID, username)
	if err != nil {
		return err
	}

	return nil
}

// GetRandomAdmin retreives random admin id from table admins
func (m *DatabaseModel) GetRandomAdmin() (int64, error) {
	stmt := `
			SELECT
				user_id
			FROM 
			    admins
			ORDER BY random()
			LIMIT 1;
		`

	var id int64
	res, err := m.DB.Query(stmt)
	if err != nil {
		return 0, err
	}

	err = res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
