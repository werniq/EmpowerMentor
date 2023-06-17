package models

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

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
	meals, nutrients, err := m.RetrieveMealsAndNutrientsByDay(userId, "monday")
	if err != nil {
		return Week{}, err
	}

	monday := Day{
		Meals:     meals,
		Nutrients: nutrients,
	}

	meals, nutrients, err = m.RetrieveMealsAndNutrientsByDay(userId, "tuesday")
	if err != nil {
		return Week{}, err
	}

	tuesday := Day{
		Meals:     meals,
		Nutrients: nutrients,
	}

	meals, nutrients, err = m.RetrieveMealsAndNutrientsByDay(userId, "wednesday")
	if err != nil {
		return Week{}, err
	}

	wednesday := Day{
		Meals:     meals,
		Nutrients: nutrients,
	}

	meals, nutrients, err = m.RetrieveMealsAndNutrientsByDay(userId, "thursday")
	if err != nil {
		return Week{}, err
	}

	thursday := Day{
		Meals:     meals,
		Nutrients: nutrients,
	}

	meals, nutrients, err = m.RetrieveMealsAndNutrientsByDay(userId, "friday")
	if err != nil {
		return Week{}, err
	}

	friday := Day{
		Meals:     meals,
		Nutrients: nutrients,
	}

	meals, nutrients, err = m.RetrieveMealsAndNutrientsByDay(userId, "saturday")
	if err != nil {
		return Week{}, err
	}

	saturday := Day{
		Meals:     meals,
		Nutrients: nutrients,
	}

	meals, nutrients, err = m.RetrieveMealsAndNutrientsByDay(userId, "sunday")
	if err != nil {
		return Week{}, err
	}

	sunday := Day{
		Meals:     meals,
		Nutrients: nutrients,
	}

	week := Week{
		Monday:    monday,
		Tuesday:   tuesday,
		Wednesday: wednesday,
		Thursday:  thursday,
		Friday:    friday,
		Saturday:  saturday,
		Sunday:    sunday,
	}

	return week, nil
}

// RetrieveMealsAndNutrientsByDay function retrieves meals and nutrients from database by day
func (m *DatabaseModel) RetrieveMealsAndNutrientsByDay(userId int64, day string) ([]Meal, Nutrients, error) {
	var mealIds []int64
	var nutrientsId int64

	stmt := fmt.Sprintf(`
			SELECT 
			    meal_ids, 
			    nutrients_id 
			FROM 
			    "%s" 
			WHERE 
			    user_id = $1
			`, day)

	err := m.DB.QueryRow(stmt, userId).Scan(pq.Array(&mealIds), &nutrientsId)
	if err != nil {
		return nil, Nutrients{}, err
	}

	var meals []Meal
	var meal Meal

	for _, id := range mealIds {
		meal, err = m.GetMealById(id)
		if err != nil {
			return nil, Nutrients{}, err
		}
		meals = append(meals, meal)
	}
	// same for nutrients
	var nutrients Nutrients
	stmt = `
			SELECT
			    calories,
			    protein,
			    fat,
			    carbohydrates
			FROM
			    nutrients
			WHERE
			    id = $1
			`

	err = m.DB.QueryRow(stmt, nutrientsId).Scan(&nutrients.Calories, &nutrients.Protein, &nutrients.Fat, &nutrients.Carbohydrates)
	if err != nil {
		return nil, Nutrients{}, err
	}

	return meals, nutrients, nil
}

// GetMealById function gets meal from database by id
func (m *DatabaseModel) GetMealById(id int64) (Meal, error) {
	var meal Meal
	stmt := `
			SELECT 
			    id, 
			    image_type, 
			    title, 
			    ready_in_minutes, 
			    servings, 
			    source_url 
			FROM 
			    meals 
			WHERE 
			    id = $1
			`

	err := m.DB.QueryRow(stmt, id).Scan(&meal.ID, &meal.ImageType, &meal.Title, &meal.ReadyInMinutes, &meal.Servings, &meal.SourceURL)
	if err != nil {
		return Meal{}, err
	}

	return meal, nil
}

// InsertWeekMealPreparePlan function inserts week meal prepare plan into database
func (m *DatabaseModel) InsertWeekMealPreparePlan(week Week, userId int64) (int64, error) {
	var ids []int64
	id, err := m.InsertMealsIntoDayTable("monday", week.Monday.Meals, week.Monday.Nutrients, userId)
	if err != nil {
		return 0, err
	}
	ids = append(ids, id)

	id, err = m.InsertMealsIntoDayTable("tuesday", week.Monday.Meals, week.Monday.Nutrients, userId)
	if err != nil {
		return 0, err
	}
	ids = append(ids, id)

	id, err = m.InsertMealsIntoDayTable("wednesday", week.Monday.Meals, week.Monday.Nutrients, userId)
	if err != nil {
		return 0, err
	}
	ids = append(ids, id)

	id, err = m.InsertMealsIntoDayTable("thursday", week.Monday.Meals, week.Monday.Nutrients, userId)
	if err != nil {
		return 0, err
	}
	ids = append(ids, id)

	id, err = m.InsertMealsIntoDayTable("friday", week.Monday.Meals, week.Monday.Nutrients, userId)
	if err != nil {
		return 0, err
	}
	ids = append(ids, id)

	id, err = m.InsertMealsIntoDayTable("saturday", week.Monday.Meals, week.Monday.Nutrients, userId)
	if err != nil {
		return 0, err
	}
	ids = append(ids, id)

	id, err = m.InsertMealsIntoDayTable("sunday", week.Monday.Meals, week.Monday.Nutrients, userId)
	if err != nil {
		return 0, err
	}
	ids = append(ids, id)

	stmt := `
		INSERT INTO
		    week(user_id, monday, tuesday, wednesday, thursday, friday, saturday, sunday)
		VALUES
		    ($1, $2, $3, $4, $5, $6, $7, $8);`

	var res sql.Result
	res, err = m.DB.Exec(stmt, userId, ids[0], ids[1], ids[2], ids[3], ids[4], ids[5], ids[6])
	if err != nil {
		return 0, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// InsertMealsIntoDayTable function inserts meals into day table. Returns id of inserted day
func (m *DatabaseModel) InsertMealsIntoDayTable(day string, meals []Meal, nutrients Nutrients, userId int64) (int64, error) {
	stmt := fmt.Sprintf(
		`
			INSERT INTO 
				meals(image_type, title, 
				 ready_in_minutes, 
				 servings, source_url) 
			VALUES 
				($1, $2, $3, $4, $5);`, day)

	var ids []int64
	var res sql.Result
	var err error
	var id int64

	for _, meal := range meals {
		res, err = m.DB.Exec(stmt, meal.ImageType, meal.Title, meal.ReadyInMinutes, meal.Servings, meal.SourceURL)
		if err != nil {
			return 0, err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return 0, err
		}
		ids = append(ids, id)
	}

	// ids -> meals_id
	// id -> nutrients_id

	stmt = "INSERT INTO nutrients (calories, protein, fat, carbohydrates) VALUES ($1, $2, $3, $4);"

	res, err = m.DB.Exec(stmt, nutrients.Calories, nutrients.Protein, nutrients.Fat, nutrients.Carbohydrates)
	if err != nil {
		return 0, err
	}

	id, err = res.LastInsertId()

	stmt = fmt.Sprintf(`INSERT INTO "%s" (user_id, meal_ids, nutrient_id) VALUES ($1, $2, $3);`, day)
	res, err = m.DB.Exec(stmt, userId, pq.Array(ids), id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

/*

	for _, meal := range week.Monday.Meals {
		res, err := m.DB.Exec(stmt, meal.ImageType, meal.Title, meal.ReadyInMinutes, meal.Servings, meal.SourceURL)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	for _, meal := range week.Tuesday.Meals {
		res, err := m.DB.Exec(stmt, meal.ImageType, meal.Title, meal.ReadyInMinutes, meal.Servings, meal.SourceURL)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	for _, meal := range week.Wednesday.Meals {
		res, err := m.DB.Exec(stmt, meal.ImageType, meal.Title, meal.ReadyInMinutes, meal.Servings, meal.SourceURL)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	for _, meal := range week.Thursday.Meals {
		res, err := m.DB.Exec(stmt, meal.ImageType, meal.Title, meal.ReadyInMinutes, meal.Servings, meal.SourceURL)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	for _, meal := range week.Friday.Meals {
		res, err := m.DB.Exec(stmt, meal.ImageType, meal.Title, meal.ReadyInMinutes, meal.Servings, meal.SourceURL)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	for _, meal := range week.Saturday.Meals {
		res, err := m.DB.Exec(stmt, meal.ImageType, meal.Title, meal.ReadyInMinutes, meal.Servings, meal.SourceURL)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	for _, meal := range week.Sunday.Meals {
		res, err := m.DB.Exec(stmt, meal.ImageType, meal.Title, meal.ReadyInMinutes, meal.Servings, meal.SourceURL)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

*/
