# Recipe Scraper

This is a recipe scraper written in Go that uses [goquery](https://github.com/PuerkitoBio/goquery) to scrap different websites and return the information written in it.

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

I might add a **generic scraper** to try and get information from any recipe site, but for now I'll work on these. Feel free to contribute!

## Structs
### Recipe
The scraper returns a `Recipe` struct in the following form:
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
Inside the [scrape](https://github.com/Stasky745/recipeScraper/tree/master/scrape) folder add a file with the domain of the site you want to scrape.

That file must contain a function like the following:

```go
func Domain(url string) []objects.Recipe {}
```

This is where all the code to extract the information from the website must go. Feel free to create any other functions necessary. If a function is not specific to your certain domain, you can add it to the [util](https://github.com/Stasky745/recipeScraper/tree/master/util) directory.

Once done this, add your domain as a `case` in [`scraper.go`](https://github.com/Stasky745/recipeScraper/blob/master/scraper.go) and create a [test](https://github.com/Stasky745/recipeScraper/tree/master/test) for it, including any particular cases that may happen in your website (eg: cases where the url contains more than one recipe). Once done this, adapt [`scraper_test.go`](https://github.com/Stasky745/recipeScraper/tree/master/test).

