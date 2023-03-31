package scrape

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"recipeScraper/objects"
	"recipeScraper/util/clean"
	"recipeScraper/util/parse"
	"regexp"
	"strconv"
	"strings"
)

func Fitmencook(url string) []objects.Recipe {
	var (
		titles      = make([]string, 0)
		description = make([]string, 0)
		keywords    = make([][]string, 0)
		steps       = make([][]objects.Steps, 1)
		ingredients = make([][]objects.Ingredient, 1)
		nutrition   = make([]objects.Nutrition, 0)
		images      = make([][]string, 0)
		servings    = make([]int, 0)
		prepTimes   = make([]int, 0)
		cookTimes   = make([]int, 0)
		video       = ""
	)

	doc := getDocFromURL(url)

	// Find the title
	doc.Find(".recipe-instructions-title").Each(func(i int, selection *goquery.Selection) {
		titles = append(titles, clean.Strings(selection.Text()))
	})

	// Find the prep-time
	doc.Find(".prep-time").Each(func(i int, selection *goquery.Selection) {
		prepTimes = append(prepTimes, time2minutes(selection.Text()))
	})

	// Find the cook-time
	doc.Find(".cook-time").Each(func(i int, selection *goquery.Selection) {
		cookTimes = append(cookTimes, time2minutes(selection.Text()))
	})

	// Find the number of servings and any images of the final dish
	doc.Find(".recipe-ingredients").Each(func(i int, selection *goquery.Selection) {
		servings = append(servings, extractServings(selection.Text()))
		images = append(images, make([]string, 0))

		selection.Find("img").Each(func(i int, selection *goquery.Selection) {
			img, _ := selection.Attr("src")
			images[len(images)-1] = append(images[len(images)-1], img)
		})
	})

	// Find the nutrition values
	doc.Find(".macros-bottom-content").Each(func(i int, selection *goquery.Selection) {
		nutrition = append(nutrition, objects.Nutrition{})

		selection.Find(".macros-bottom-info").Each(func(i int, selection *goquery.Selection) {
			name := strings.TrimSpace(selection.Find(".macros-label").First().Text())
			value := strings.TrimSpace(selection.Find(".macros").First().Text())
			nut := parse.NutritionValue(value)

			switch strings.ToLower(name) {
			case "calories":
				nutrition[len(nutrition)-1].Calories = nut
			case "protein":
				nutrition[len(nutrition)-1].Protein = nut
			case "fat":
				nutrition[len(nutrition)-1].Fat = nut
			case "carbs":
				nutrition[len(nutrition)-1].Carbs = nut
			case "sodium":
				nutrition[len(nutrition)-1].Sodium = nut
			case "fiber":
				nutrition[len(nutrition)-1].Fiber = nut
			case "sugar":
				nutrition[len(nutrition)-1].Sugar = nut
			}
		})
	})

	// Find the ingredients list
	// If there are multiple recipes in a single page, we determine which recipe we're getting the ingredients from by comparing the parents, increasing the counter whenever it changes
	ingredientCounter := 0
	parentSelection := ""
	doc.Find(".recipe-ingredients > ul > li").Each(func(i int, selection *goquery.Selection) {
		if parentSelection == "" {
			parentSelection = selection.Parent().First().Text()
		} else {
			newParentSelection := selection.Parent().First().Text()
			if newParentSelection != parentSelection {
				parentSelection = newParentSelection
				ingredientCounter++
				ingredients = append(ingredients, make([]objects.Ingredient, 0))
			}
		}

		ingredients[ingredientCounter] = append(ingredients[ingredientCounter], getIngredients(selection))
	})

	stepsCounter := -1
	parentSelection = ""
	doc.Find(".recipe-steps").Each(func(i int, selection *goquery.Selection) {
		stepsCounter++
		steps = append(steps, make([]objects.Steps, 0))
		step := 1

		selection.Find("p").Each(func(i int, selection *goquery.Selection) {
			if steps[stepsCounter] == nil || len(steps[stepsCounter]) == 0 {
				steps[stepsCounter] = append(steps[stepsCounter], objects.Steps{Step: step})
				step++
			}

			lastStep := &steps[stepsCounter][len(steps[stepsCounter])-1]
			if text := selection.Text(); text != "" {
				if lastStep.Instruction == "" {
					lastStep.Instruction = text
				} else {
					newStep := objects.Steps{
						Instruction: text,
						Picture:     "",
						Step:        step,
					}
					step++
					steps[stepsCounter] = append(steps[stepsCounter], newStep)
				}
			} else {
				img, _ := selection.Find("img").First().Attr("src")
				if lastStep.Picture == "" {
					lastStep.Picture = img
				} else {
					newStep := objects.Steps{
						Instruction: "",
						Picture:     img,
						Step:        step,
					}
					step++
					steps[stepsCounter] = append(steps[stepsCounter], newStep)
				}
			}
		})
	})

	// Videos
	video = doc.Find(".contains-sticky-video iframe").First().AttrOr("src", "")

	// There are no descriptions/keywords in this site, so we create an empty array
	description = make([]string, len(titles))
	keywords = make([][]string, len(titles))

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
			Images:      images[i],
			Servings:    servings[i],
			PrepTime:    prepTimes[i],
			CookTime:    cookTimes[i],
			Video:       video,
		}

		recipes = append(recipes, recipe)
	}
	return recipes
}

// Returns the ingredients for a recipe. Uses recursion in case an ingredient has sub-ingredients
func getIngredients(selection *goquery.Selection) objects.Ingredient {
	sublists := selection.Find("ul")
	parent := sublists.Parent().Find("strong").First().Text()

	var subIngredient objects.Ingredient
	if sublists.Length() > 0 {
		subIngredient = objects.Ingredient{
			Name:        parent,
			Quantity:    1,
			Unit:        "",
			Ingredients: make([]objects.Ingredient, 0),
			Notes:       "",
		}

		sublists.Find("li").Each(func(i int, item *goquery.Selection) {
			subparent := item.Parent().Parent().Find("strong").First().Text()
			if parent == subparent {
				subIngredient.Ingredients = append(subIngredient.Ingredients, getIngredients(item))
			}
		})
	} else {
		subIngredient = parse.Ingredient(selection.Text())
	}

	return subIngredient
}

func extractServings(s string) int {
	reg := regexp.MustCompile(`.*ngredients for (\d+) servings.*`)
	string := reg.FindStringSubmatch(s)

	servings, err := strconv.Atoi(string[1])
	if err != nil {
		log.Fatalln("Error at converting Servings from string to int")
	}

	return servings
}

// Gets a string in the format of "X hours Y minutes" and returns it as an int in minutes
func time2minutes(s string) int {
	reg := regexp.MustCompile(`(?:(?P<hours>\d+)\s*[Hh]ours?\s*)?(?:(?P<mins>\d+)\s*[Mm]inutes?)?`)
	match := reg.FindStringSubmatch(s)

	result := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	minutes, err := strconv.Atoi(result["mins"])
	if err != nil {
		minutes = 0
	}

	hours, err := strconv.Atoi(result["hours"])
	if err != nil {
		hours = 0
	}

	return hours*60 + minutes
}
