package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"self-improvement-bot/models"
)

type ConnectBody struct {
	Username  string `json:"username"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastN ame"`
	Email     string `json:"email"`
}

type SpoonacularResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hash     string `json:"hash"`
}

type SpoonacularConfig struct {
	Config SpoonacularResponse `json:"config"`
}

type MealPlan struct {
	Days []Day `json:"days"`
}
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

var (
	SpoonacularURI = "https://api.spoonacular.com/users/connect?apiKey=" + os.Getenv("SPOONOCULAR_API_KEY")

	// MealPrepareRequestUri = "https://api.spoonacular.com/mealplanner/qniwerniq1/week/2023-06-06?hash=bda3c8447f862d6927612543fbe3f8bef3af03c2&apiKey=7bd3812e64cc4b3bbde5ac72cd575331"

	MealPrepareGenerator = "https://api.spoonacular.com/mealplanner/generate"

	SpoonocularTypeOfDiest = []string{}
)

// ConnectSpoonacular connects to Spoonacular API` and returns a response with username, password and hash
func ConnectSpoonacular() (*SpoonacularConfig, error) {
	var ConnectBody ConnectBody
	ConnectBody.Firstname = "Oleksandr"
	ConnectBody.Lastname = "Matviienko"
	ConnectBody.Email = "qniwwwersss@gmail.com"
	ConnectBody.Username = "qniwerniq"

	req, err := http.NewRequest("POST", SpoonacularURI, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	response := SpoonacularResponse{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	q := &SpoonacularConfig{Config: response}
	return q, nil
}

// CreateMealPreparingPlan creates a meal preparing plan
func CreateMealPreparingPlan(uri string) (models.Week, error) {
	//
	res, err := http.Get(uri)
	if err != nil {
		return models.Week{}, err
	}

	var week struct {
		Week models.Week `json:"week"`
	}
	err = json.NewDecoder(res.Body).Decode(&week)
	if err != nil {
		return models.Week{}, err
	}
	fmt.Println(week)
	fmt.Println(week)
	fmt.Println(week)

	return week.Week, nil
}
