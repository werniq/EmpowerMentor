package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	Source      Source    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

type NewsApiResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

var (
	NewsApiUri          = "https://newsapi.org/v2/everything?q="
	AuthenticateUriPart = "&apiKey=6082fbd18c054e72a0077c422436fdf2"

	NewsApiCategories = []string{
		"health",
		"fitness",
		"personal",
		"growth",
		"psychology",
		"mindfulness",
		"self-Care",
		"motivation",
		"productivity",
		"happiness",
		"relationships",
		"career-development",
		"leadership",
		"education",
		"self-Help",
		"mental",
		"well-being",
	}
)

func RetrieveCategories() []string {
	return AvaliableCategories
}

// RequestToNewsApi makes a request to newsapi.org and returns a NewsApiResponse
func RequestToNewsApi(category string) (*NewsApiResponse, error) {
	res := &NewsApiResponse{}

	ok := false
	for i := 0; i <= len(NewsApiCategories)-1; i++ {
		if category == NewsApiCategories[i] {
			ok = true
		}
	}

	if !ok {
		return nil, errors.New("category not found")
	}

	url := NewsApiUri + category + AuthenticateUriPart
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
