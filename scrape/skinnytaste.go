package scrape

import (
	"github.com/PuerkitoBio/goquery"
	"recipeScraper/objects"
	"recipeScraper/util/clean"
	"recipeScraper/util/parse"
	"strconv"
)

func Skinnytaste(url string) []objects.Recipe {
	var (
		titles      = make([]string, 0)
		description = make([]string, 0)
		keywords    = make([][]string, 0)
		steps       = make([][]objects.Steps, 0)
		ingredients = make([][]objects.Ingredient, 0)
		nutrition   = make([]objects.Nutrition, 0)
		images      = make([]string, 0)
		servings    = make([]int, 0)
		prepTimes   = make([]int, 0)
		cookTimes   = make([]int, 0)
		videos      = make([]string, 0)
	)

	doc := getDocFromURL(url)

	// Find the title
	doc.Find(".wprm-recipe-name").Each(func(i int, selection *goquery.Selection) {
		titles = append(titles, clean.Strings(selection.Text()))
	})

	// Find the description
	doc.Find(".wprm-recipe-summary").Each(func(i int, selection *goquery.Selection) {
		description = append(description, clean.Strings(selection.Text()))
	})

	// Find the prep time
	doc.Find(".wprm-recipe-prep-time-container").Each(func(i int, selection *goquery.Selection) {
		prepTime := 0

		//Find the hours
		doc.Find(".wprm-recipe-prep_time-hours").Each(func(i int, selection *goquery.Selection) {
			hours, err := strconv.Atoi(selection.Text())
			if err != nil {
				hours = 0
			}

			prepTime += hours * 60
		})

		//Find the minutes
		doc.Find(".wprm-recipe-prep_time-minutes").Each(func(i int, selection *goquery.Selection) {
			minutes, err := strconv.Atoi(selection.Text())
			if err != nil {
				minutes = 0
			}

			prepTime += minutes
		})

		prepTimes = append(prepTimes, prepTime)
	})

	// Find the cook time
	doc.Find(".wprm-recipe-cook-time-container").Each(func(i int, selection *goquery.Selection) {
		cookTime := 0

		//Find the hours
		doc.Find(".wprm-recipe-cook_time-hours").Each(func(i int, selection *goquery.Selection) {
			hours, err := strconv.Atoi(selection.Text())
			if err != nil {
				hours = 0
			}

			cookTime += hours * 60
		})

		//Find the minutes
		doc.Find(".wprm-recipe-cook_time-minutes").Each(func(i int, selection *goquery.Selection) {
			minutes, err := strconv.Atoi(selection.Text())
			if err != nil {
				minutes = 0
			}

			cookTime += minutes
		})

		cookTimes = append(cookTimes, cookTime)
	})

	// Find the servings
	doc.Find(".wprm-recipe-servings").Each(func(i int, selection *goquery.Selection) {
		serv, err := strconv.Atoi(selection.Text())
		if err != nil {
			serv = 1
		}

		servings = append(servings, serv)
	})

	// Find the ingredients
	doc.Find(".wprm-recipe-ingredients-container").Each(func(i int, selection *goquery.Selection) {
		newIngredients := make([]objects.Ingredient, 0)

		selection.Find(".wprm-recipe-ingredient-group").Each(func(i int, selection *goquery.Selection) {
			subIngredients := make([]objects.Ingredient, 0)
			groupName := selection.Find(".wprm-recipe-ingredient-group-name").First().Text()

			selection.Find(".wprm-recipe-ingredient").Each(func(i int, selection *goquery.Selection) {
				subIngredients = append(subIngredients, parse.Ingredient(selection.Text()))
			})

			if groupName != "" {
				ingredient := objects.Ingredient{
					Name:        groupName,
					Quantity:    1,
					Unit:        "",
					Ingredients: subIngredients,
					Notes:       "",
				}

				newIngredients = append(newIngredients, ingredient)
			} else {
				newIngredients = append(newIngredients, subIngredients...)
			}
		})

		ingredients = append(ingredients, newIngredients)
	})

	// Instructions
	doc.Find(".wprm-recipe-instruction-group").Each(func(i int, selection *goquery.Selection) {
		substeps := make([]objects.Steps, 0)
		selection.Find(".wprm-recipe-instruction-text").Each(func(i int, selection *goquery.Selection) {
			step := objects.Steps{
				Step:        i + 1,
				Instruction: selection.Text(),
				Picture:     "",
			}

			substeps = append(substeps, step)
		})

		steps = append(steps, substeps)
	})

	// Find nutrition values
	doc.Find(".wprm-nutrition-label-container").Each(func(i int, selection *goquery.Selection) {
		nut := objects.Nutrition{
			Calories: extractNutritionValues(".wprm-nutrition-label-text-nutrition-container-calories", selection),
			Protein:  extractNutritionValues(".wprm-nutrition-label-text-nutrition-container-protein", selection),
			Fat:      extractNutritionValues(".wprm-nutrition-label-text-nutrition-container-fat", selection),
			Carbs:    extractNutritionValues(".wprm-nutrition-label-text-nutrition-container-carbohydrates", selection),
			Sodium:   extractNutritionValues(".wprm-nutrition-label-text-nutrition-container-sodium", selection),
			Fiber:    extractNutritionValues(".wprm-nutrition-label-text-nutrition-container-fiber", selection),
			Sugar:    extractNutritionValues(".wprm-nutrition-label-text-nutrition-container-sugar", selection),
		}

		nutrition = append(nutrition, nut)
	})

	// Images
	doc.Find(".site-main .entry-content").Each(func(i int, selection *goquery.Selection) {
		imageMap := make(map[string]struct{}, 0)
		selection.Find("img").Each(func(i int, selection *goquery.Selection) {
			imgLink, exists := selection.Attr("data-pin-media")
			if exists {
				imageMap[imgLink] = struct{}{}
			}
		})

		if len(imageMap) > 0 {
			for k := range imageMap {
				images = append(images, k)
			}
		}
	})

	// Video
	videos = append(videos, doc.Find(".wprm-recipe-video .rll-youtube-player").First().AttrOr("data-src", ""))

	// There are no keywords in this site
	keywords = make([][]string, len(titles))

	if len(prepTimes) < len(titles) {
		prepTimes = make([]int, len(titles))
	}

	if len(cookTimes) < len(titles) {
		cookTimes = make([]int, len(titles))
	}

	recipes := make([]objects.Recipe, 0)
	for i := 0; i < len(titles); i++ {
		recipe := objects.Recipe{
			Title:       titles[i],
			Description: description[i],
			Url:         url,
			Keywords:    keywords[i],
			Steps:       steps[i],
			Ingredients: ingredients[i],
			Nutrition:   nutrition[i],
			Images:      images,
			Servings:    servings[i],
			PrepTime:    prepTimes[i],
			CookTime:    cookTimes[i],
			Video:       videos[i],
		}

		recipes = append(recipes, recipe)
	}
	return recipes
}

func extractNutritionValues(class string, selection *goquery.Selection) objects.NutritionValue {
	text := selection.Find(class + " > .wprm-nutrition-label-text-nutrition-value").First().Text()
	val, err := strconv.Atoi(text)
	if err != nil {
		val = 1
	}
	unit := selection.Find(class + " > .wprm-nutrition-label-text-nutrition-unit").First().Text()

	return objects.NutritionValue{
		Value: val,
		Unit:  unit,
	}
}
