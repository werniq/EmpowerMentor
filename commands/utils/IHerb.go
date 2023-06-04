package utils

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

type IherbItem struct {
	Title string
	Price string
}

var (
	AvaliableCategories = []string{
		"vitamins",
		"supplements",
		"minerals",
		"digestive-support",
		"antioxidants",
		"bone-joint-cartilage",
		"sleep",
		"fish-oil-omegas-epa-dha",
		"brain-cognitive",
		"hair-skin-nails",
		"greens-superfoods",
		"amino-acids",
		"bee-products",
		"childrens-health",
		"mens-health",
		"mushrooms",
		"weight-loss",
		"phospholipids",
		"protein",
		"womens-health",
	}
)

// VerifyCategory verifies if category is available
func VerifyCategory(category string) bool {
	for i := 0; i <= len(AvaliableCategories)-1; i++ {
		if AvaliableCategories[i] == category {
			return true
		}
	}
	return false
}

// RequestToIHerb makes a request to iherb.com and returns a slice of IherbItem
func RequestToIHerb(category string) ([]IherbItem, error) {
	var ok bool
	for i := 0; i <= len(AvaliableCategories)-1; i++ {
		if category == AvaliableCategories[i] {
			ok = true
		}
	}
	if !ok {
		return nil, fmt.Errorf("category not found. avaliable categories: %v", AvaliableCategories)
	}

	url := "https://pl.iherb.com/c/" + category

	res, err := http.Get(url)
	if err != nil {
		log.Printf("error getting response: %v\n", err)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("error creating document from reader: %v\n", err)
		return nil, err
	}

	var items []IherbItem

	for a := 0; a < 10; a++ {
		var item IherbItem
		_ = doc.WrapAll(".product-title").Each(func(i int, s *goquery.Selection) {
			item.Title = s.Find(".product-title").Text()
			item.Price = s.Find(".price").Text()
		})

		items = append(items, item)
	}

	return items, nil
}
