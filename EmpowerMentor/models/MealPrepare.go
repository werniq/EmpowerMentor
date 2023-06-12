package models

type Meal struct {
	ID             int    `json:"id"`
	ImageType      string `json:"imageType"`
	Title          string `json:"title"`
	ReadyInMinutes int    `json:"readyInMinutes"`
	Servings       int    `json:"servings"`
	SourceURL      string `json:"sourceUrl"`
}

type Nutrients struct {
	Calories      float64 `json:"calories"`
	Protein       float64 `json:"protein"`
	Fat           float64 `json:"fat"`
	Carbohydrates float64 `json:"carbohydrates"`
}

type Day struct {
	Meals     []Meal    `json:"meals"`
	Nutrients Nutrients `json:"nutrients"`
}

type Week struct {
	Monday    Day `json:"monday"`
	Tuesday   Day `json:"tuesday"`
	Wednesday Day `json:"wednesday"`
	Thursday  Day `json:"thursday"`
	Friday    Day `json:"friday"`
	Saturday  Day `json:"saturday"`
	Sunday    Day `json:"sunday"`
}

// GetMealPlan function gets meal plan from database
func (m *DatabaseModel) GetMealPlan(userId int64) (Week, error) {
	stmt := `
			SELECT 
				monday, tuesday, wednesday, thursday, friday, saturday, sunday
			FROM
			    week
			WHERE
			    user_id = $1`

	var week Week
	err := m.DB.QueryRow(stmt, userId).Scan(&week.Monday, &week.Tuesday, &week.Wednesday, &week.Thursday, &week.Friday, &week.Saturday, &week.Sunday)
	if err != nil {
		return Week{}, err
	}

	return week, nil
}

// InsertMealPreparePlan function inserts meal prepare plan into database
func (m *DatabaseModel) InsertMealPreparePlan(week Week, userId int64) error {
	stmt := `
		INSERT INTO
		    week(user_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err := m.DB.Exec(stmt, userId, week.Monday, week.Tuesday, week.Wednesday, week.Thursday, week.Friday, week.Saturday, week.Sunday)
	if err != nil {
		return err
	}

	return nil
}
