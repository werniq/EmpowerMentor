package models

import (
	"database/sql"
	"strings"
	"time"
)

type UserBotConfiguration struct {
	UserId                    int64                `json:"user_id"`
	Username                  string               `json:"username"`
	Step                      int                  `json:"step"`
	Gender                    string               `json:"gender"`
	Age                       int                  `json:"age"`
	Weight                    float64              `json:"weight"`
	Height                    float64              `json:"height"`
	PreferredPhysicalActivity string               `json:"preferred_physical_activity"`
	CustomReminders           map[string]time.Time `json:"custom_reminders"`
	WorkoutCount              int                  `json:"workout_count"`
	BooksCount                int                  `json:"books_count"`
	PreferringSupplements     string               `json:"preferring_supplements"`
	HabitsToAcquire           string               `json:"habits_to_acquire"`
	NewsCategories            string               `json:"news_categories"`
	WakeUpTime                time.Time            `json:"wake_up_time"`
	BedTime                   time.Time            `json:"bed_time"`
	PreferableTimeToMeditate  time.Time            `json:"preferable_time_to_meditate"`
	PreferableTimeToExercise  time.Time            `json:"preferable_time_to_exercise"`
	PreferableTimeToRead      time.Time            `json:"preferable_time_to_read"`
}

type Motivation struct {
	UserID     int64  `json:"user_id"`
	Motivation string `json:"motivation"`
}

type DatabaseModel struct {
	DB *sql.DB
}

// NewDatabaseModel creates a new DatabaseModel struct
func NewDatabaseModel(db *sql.DB) *DatabaseModel {
	return &DatabaseModel{DB: db}
}

// GetDailyMeditationReminderForUser returns time when bot sends daily meditation reminder to user
func (m *DatabaseModel) GetDailyMeditationReminderForUser(userId int64) (time.Time, error) {
	stmt := "SELECT preferable_time_to_meditate FROM user_bot_configuration WHERE user_id = $1"
	var result time.Time
	res := m.DB.QueryRow(stmt, userId)
	if err := res.Scan(&result); err != nil {
		return time.Time{}, err
	}
	return result, nil
}

// GetDailyExerciseReminderForUser returns time when bot sends daily exercise reminder to user
func (m *DatabaseModel) GetDailyExerciseReminderForUser(userId int64) (time.Time, error) {
	stmt := "SELECT preferable_time_to_exercise FROM user_bot_configuration WHERE user_id = $1"

	var result time.Time

	row, err := m.DB.Query(stmt, userId)
	if err != nil {
		return time.Time{}, err
	}

	err = row.Scan(&result)
	if err != nil {
		return time.Time{}, err
	}

	return result, nil
}

// GetDailyReadingReminderForUser returns time when bot sends daily reading reminder to user
func (m *DatabaseModel) GetDailyReadingReminderForUser(userId int64) (time.Time, error) {
	stmt := "SELECT preferable_time_to_read FROM user_bot_configuration WHERE user_id = $1"
	var result time.Time

	row, err := m.DB.Query(stmt, userId)
	if err != nil {
		return time.Time{}, err
	}

	err = row.Scan(&result)
	if err != nil {
		return time.Time{}, err
	}

	return result, nil
}

// SetDailyMeditationReminderForUser resets time when bot sends daily meditation reminder to user
func (m *DatabaseModel) SetDailyMeditationReminderForUser(userId int64, time time.Time) error {
	stmt := "UPDATE user_bot_configuration SET preferable_time_to_meditate = $1 WHERE user_id = $2"
	_, err := m.DB.Exec(stmt, time, userId)
	if err != nil {
		return err
	}
	return nil
}

// SetDailyExerciseReminderForUser resets time when bot sends daily exercise reminder to user
func (m *DatabaseModel) SetDailyExerciseReminderForUser(userId int64, time time.Time) error {
	stmt := "UPDATE user_bot_configuration SET preferable_time_to_exercise = $1 WHERE user_id = $2"
	_, err := m.DB.Exec(stmt, time, userId)
	if err != nil {
		return err
	}
	return nil
}

// SetDailyReadingReminderForUser resets time when bot sends daily reading reminder to user
func (m *DatabaseModel) SetDailyReadingReminderForUser(userId int64, time time.Time) error {
	stmt := "UPDATE user_bot_configuration SET preferable_time_to_read = $1 WHERE user_id = $2"
	_, err := m.DB.Exec(stmt, time, userId)
	if err != nil {
		return err
	}
	return nil
}

// SetCustomReminderForUser resets time when bot sends daily reminder for doing something to user
func (m *DatabaseModel) SetCustomReminderForUser(userId int64, reminder string) (time.Time, error) {
	stmt := "SELECT preferable_time_to_" + reminder + " FROM user_bot_configuration WHERE user_id = $1"

	var result time.Time

	row, err := m.DB.Query(stmt, userId)
	if err != nil {
		return time.Time{}, err
	}

	if err = row.Scan(&result); err != nil {
		return time.Time{}, err
	}

	return result, nil
}

// DeleteReminder deletes reminder for user
func (m *DatabaseModel) DeleteReminder(userId int64, reminder string) error {
	stmt := `
			UPDATE 
			    user_bot_configuration 
			SET 
			    preferable_time_to_$1 = NULL 
			WHERE 
			    user_id = $2`

	_, err := m.DB.Exec(stmt, reminder, userId)
	if err != nil {
		return err
	}

	return nil
}

// GetUserReminders returns all reminders for user
func (m *DatabaseModel) GetUserReminders(userId int64) ([]time.Time, []string, error) {
	stmt := "SELECT preferable_time_to_meditate, preferable_time_to_exercise, preferable_time_to_read FROM user_bot_configuration WHERE user_id = $1"

	var meditation, exercise, reading time.Time

	row, err := m.DB.Query(stmt, userId)
	if err != nil {
		return nil, nil, err
	}

	if err = row.Scan(&meditation, &exercise, &reading); err != nil {
		return nil, nil, err
	}

	var times []time.Time
	times[0] = meditation
	times[1] = exercise
	times[2] = reading

	return times, []string{"meditate, exercise, read"}, nil
}

// DailyMeditationReminderForUser returns true if it is time to send daily meditation reminder to user
func (m *DatabaseModel) DailyMeditationReminderForUser(userId int64) (bool, error) {
	t, err := m.GetDailyMeditationReminderForUser(userId)
	if err != nil {
		return false, err
	}

	if t.After(time.Now()) {
		return true, nil
	}

	return false, nil
}

// DailyExerciseReminderForUser returns true if it is time to send daily exercise reminder to user
func (m *DatabaseModel) DailyExerciseReminderForUser(userId int64) (bool, error) {
	t, err := m.GetDailyExerciseReminderForUser(userId)
	if err != nil {
		return false, err
	}

	if t.After(time.Now()) {
		return true, nil
	}

	return false, nil
}

//func (m *DatabaseModel) DailyHabitReminderForUser(userId int64) (bool, error)      {
//	t, err := m.GetDailyHabitReminderForUser(int64)
//}

// DailyReadingReminderForUser returns true if it is time to send daily reading reminder to user
func (m *DatabaseModel) DailyReadingReminderForUser(userId int64) (bool, error) {
	t, err := m.GetDailyReadingReminderForUser(userId)
	if err != nil {
		return false, err
	}

	if t.After(time.Now()) {
		return true, nil
	}

	return false, nil
}

// AddHabit
func (m *DatabaseModel) AddHabit(userId int64, habit string) error {
	stmt := `
		SELECT 
		    habits_to_acquire
		FROM 
		    user_bot_configuration
		WHERE 
		    user_id = $1
				`

	var habits string

	row, err := m.DB.Query(stmt, userId)
	if err != nil {
		return err
	}

	if err = row.Scan(&habits); err != nil {
		return err
	}

	habits += habit

	stmt = `
			UPDATE
			    user_bot_configuration
			SET
			    habits_to_acquire = $1
			WHERE
			    user_id = $2
				`

	return nil
}

// DeleteHabit deletes habit for user
func (m *DatabaseModel) DeleteHabit(userId int64) error {
	stmt := `
		update user_bot_configuration
		set habits_to_acquire = $1
		where user_id = $2
		`

	_, err := m.DB.Exec(stmt, "", userId)
	if err != nil {
		return err
	}

	return nil
}

// UpdateHabit updates habit for user
func (m *DatabaseModel) UpdateHabit(userId int64, habit string) error {
	h, err := m.GetHabits(userId)
	habits := strings.Split(h, " ")

	if err != nil {
		return err
	}

	stmt := "UPDATE user_bot_configuration SET habits_to_acquire = $1 WHERE user_id = $2"

	habits = append(habits, habit)
	_, err = m.DB.Exec(stmt, habits, userId)
	if err != nil {
		return err
	}

	return nil
}

// GetHabits retrieves habits for user
func (m *DatabaseModel) GetHabits(userId int64) (string, error) {
	stmt := `
		SELECT 
			habits_to_acquire
		FROM
		    user_bot_configuration
		WHERE
		    user_id = $1
		`

	var habits string
	err := m.DB.QueryRow(stmt, userId).Scan(&habits)
	if err != nil {
		return "", err
	}

	return habits, nil
}

// StoreUserData stores user data in database
func (m *DatabaseModel) StoreUserData(u UserBotConfiguration) error {
	stmt := `
			INSERT INTO 
				user_bot_configuration(
				    user_id,
					username,
					gender,
					age,
					weight,
					height,
					preferred_physical_activity,
					custom_reminders,
					workout_count,
					books_count,
					preferring_supplements,
					habits_to_acquire,
					news_categories,
					wake_up_time,
					bed_time,
					preferable_time_to_meditate,
					preferable_time_to_exercise,
					preferable_time_to_read)
			VALUES (
			        $1, $2, $3, 
			        $4, $5, $6, 
			        $7, $8, $9, 
			        $10, $11, $12, 
			        $13, $14, $15, 
			        $16, $17);`

	_, err := m.DB.Exec(stmt,
		u.UserId,
		u.Username,
		u.Gender,
		u.Age,
		u.Weight,
		u.Height,
		u.PreferredPhysicalActivity,
		u.CustomReminders,
		u.WorkoutCount,
		u.BooksCount,
		u.PreferringSupplements,
		u.HabitsToAcquire,
		u.NewsCategories,
		u.WakeUpTime,
		u.BedTime,
		u.PreferableTimeToMeditate,
		u.PreferableTimeToExercise,
		u.PreferableTimeToRead)

	if err != nil {
		return err
	}

	return nil
}

// SaveUserCustomMotivation saves user custom motivation which will be sent one time per day
func (m *DatabaseModel) SaveUserCustomMotivation(userId int64, motivation string) error {
	stmt := `
			INSERT INTO
				users_motivation(user_id, motivation)
			VALUES ($1, $2)
		`

	_, err := m.DB.Exec(stmt, userId, motivation)
	if err != nil {
		return err
	}

	return nil
}

// GetUserCustomMotivation retrieves custom motivation from
func (m *DatabaseModel) GetUserCustomMotivation(userid int64) (string, error) {
	stmt := `
		SELECT 
		    motivation 
		FROM 
		    users_motivation
		WHERE 
		    user_id = $1;
		    `

	var motivation string
	err := m.DB.QueryRow(stmt, userid).Scan(&motivation)
	if err != nil {
		return "", err
	}

	return motivation, nil
}

// CreateCustomGoal creates new record in table user_goals
func (m *DatabaseModel) CreateCustomGoal(userId int64, goal string) error {
	stmt := `
			INSERT INTO 
			    user_goals
			    (user_id, goal, created_at) 
			VALUES 
			    ($1, $2, $3)`

	_, err := m.DB.Exec(stmt, userId, goal, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// CreateCustomReminder creates new record in table user_reminders
func (m *DatabaseModel) CreateCustomReminder(userId int64, reminder string, appearTime time.Time) error {
	stmt := `
		INSERT INTO user_reminders
			(user_id, reminder, appear_time)
		VALUES 
		    ($1, $2, $3)
		`

	t := appearTime.Format("15:15:15")
	_, err := m.DB.Exec(stmt, userId, reminder, t)
	if err != nil {
		return err
	}

	return nil
}

// SaveUserCustomChallenge saves user's custom challenge
func (m *DatabaseModel) SaveUserCustomChallenge(userId int64, challenge string) error {
	stmt := `
		INSERT INTO user_challenges
			(user_id, challenge)
		VALUES
			($1, $2)
		`

	_, err := m.DB.Exec(stmt, userId, challenge)
	if err != nil {
		return err
	}

	return nil
}

//func (m *DatabaseModel) SetDailyHabitReminderForUser(int64 string) error       {
//
//	stmt := "SELECT preferable_time_to" + "_mediate FROM user_bot_configuration WHERE user_id = $1"
//}

// DailyCustomReminderForUser returns true if it is time to send daily reminder for doing something to user
//func (m *DatabaseModel) DailyCustomReminderForUser(int64 string) (bool, error)     {
//	stmt := ""
//}
//func (m *DatabaseModel) UserVerifyReminderExists(userId int64) error {}
