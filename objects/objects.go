package objects

import "recipeScraper/util/equal"

type NutritionValue struct {
	Unit  string
	Value int
}

type Nutrition struct {
	Calories NutritionValue
	Protein  NutritionValue
	Fat      NutritionValue
	Carbs    NutritionValue
	Sodium   NutritionValue
	Fiber    NutritionValue
	Sugar    NutritionValue
}

type Ingredient struct {
	Name        string
	Quantity    float64
	Unit        string
	Ingredients []Ingredient
	Notes       string
}

type Steps struct {
	Step        int
	Instruction string
	Picture     string
}

type Recipe struct {
	Title       string
	Description string
	Url         string
	Keywords    []string
	Steps       []Steps
	Ingredients []Ingredient
	Nutrition   Nutrition
	Images      []string
	Servings    int
	PrepTime    int
	CookTime    int
	Video       string
}

func (i Ingredient) Equals(i2 Ingredient) bool {
	return i.Name == i2.Name &&
		i.Unit == i2.Unit &&
		i.Quantity == i2.Quantity &&
		i.Notes == i2.Notes &&
		compareIngredientSlices(i.Ingredients, i2.Ingredients)
}

func compareIngredientSlices(i []Ingredient, i2 []Ingredient) bool {
	if len(i) != len(i2) {
		return false
	}

	for n := range i {
		if !i[n].Equals(i2[n]) {
			return false
		}
	}
	return true
}

func (r Recipe) Equals(r2 Recipe) bool {
	return r.Title == r2.Title &&
		r.Description == r2.Description &&
		r.Url == r2.Url &&
		equal.Slices(r.Keywords, r2.Keywords) &&
		equal.Slices(r.Steps, r2.Steps) &&
		compareIngredientSlices(r.Ingredients, r2.Ingredients) &&
		r.Nutrition == r2.Nutrition &&
		equal.Slices(r.Images, r2.Images) &&
		r.Servings == r2.Servings &&
		r.PrepTime == r2.PrepTime &&
		r.CookTime == r2.CookTime
}
