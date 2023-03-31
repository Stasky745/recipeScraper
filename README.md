# Recipe Scrapper

This is a recipe scrapper written in Go that uses [goquery](https://github.com/PuerkitoBio/goquery) to scrap different websites and return the information written in it.

## Sites
The sites currently working are:
1. [Skinny Taste](https://www.skinnytaste.com)
2. [Fit Men Cook](https://fitmencook.com)

This is a work in progress. Some of the websites I plan on including in no particular order in the future are:
1. [RecipeTinEats](https://www.recipetineats.com/)
2. [Blue Apron](https://www.blueapron.com)
3. [Hello Fresh](https://www.hellofresh.com/)
4. [Budget Bytes](https://www.budgetbytes.com/)
5. [Home Chef](https://www.homechef.com)

I might add a generic scrapper to try and get information from any recipe site, but for now I'll work on these. Feel free to contribute!

## Structs
### Recipe
The scrapper returns a `Recipe` struct in the following form:
```go
type Recipe struct {
	// Name of the recipe
	Title       string
	// Description of the recipe
	Description string
	// Url
	Url         string
	// Any keywords (eg: vegan, fried, etc.)
	Keywords    []string
	// Instructions to follow
	Steps       []Steps
	// Ingredients included
	Ingredients []Ingredient
	// Nutrition values (eg: calories, fat, etc.)
	Nutrition   Nutrition
	// Urls of any images of the recipe
	Images      []string
	// The number of servings this recipe is for
	Servings    int
	// Time it takes to prep in minutes
	PrepTime    int
	// Time it takes to cook in minutes
	CookTime    int
	// Video url of the recipe
	Video       string
}
```

### Steps
```go
type Steps struct {
	// Number of step
	Step        int
	// Instruction to follow
	Instruction string
	// Url of picture
	Picture     string
}
```

### Ingredient
The `Ingredient` struct can contain other ingredients in case these ingredients are for making another ingredient in the recipe (eg: a sauce).
```go
type Ingredient struct {
	// Name of the ingredient
	Name        string
	// Quantity
	Quantity    float64
	// Unit of the quantity
	Unit        string
	// Any subingredients if any
	Ingredients []Ingredient
	// Notes
	Notes       string
}
```

In the case where there are subingredients, the `Quantity` and `Unit` won't have a value.

### Nutrition
```go
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
```

## Want to contribute?
