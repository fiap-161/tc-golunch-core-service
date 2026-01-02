package enum

import "strings"

type Category string

const (
	Meal    Category = "MEAL"
	Side    Category = "SIDE"
	Drink   Category = "DRINK"
	Dessert Category = "DESSERT"
)

func IsValidCategory(category string) bool {
	category = strings.ToUpper(category)
	switch Category(category) {
	case Meal, Side, Drink, Dessert:
		return true
	default:
		return false
	}
}

func GetAllCategories() []Category {
	return []Category{Meal, Side, Drink, Dessert}
}
