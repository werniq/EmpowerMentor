package models

import (
	"time"
)

type UserConfiguration struct {
	UserId int64 `json:"user_id"`
	Step   int   `json:"step"`

	WorkoutCount             int       `json:"workout_count"`
	BooksCount               int       `json:"books_count"`
	NewsCategories           string    `json:"news_categories"`
	WakeUpTime               time.Time `json:"wake_up_time"`
	BedTime                  time.Time `json:"bed_time"`
	PreferableTimeToMeditate time.Time `json:"preferable_time_to_meditate"`
	PreferableTimeToExercise time.Time `json:"preferable_time_to_exercise"`
	PrefrableTimeToRead      time.Time `json:"preferable_time_to_read"`
}

// IsUserRegistered checks if user exists in the database
func (m *DatabaseModel) IsUserRegistered(userId int64) (bool, error) {
	stmt := `
			SELECT 
				COUNT(*)
			FROM
			    users
			WHERE
			    id = $1`

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

// RegisterUser registers user in the database
func (m *DatabaseModel) RegisterUser(userid int64) error {

	return nil
}
