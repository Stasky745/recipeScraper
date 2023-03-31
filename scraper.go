package recipeScraper

import (
	"log"
	"net/url"
	"recipeScraper/objects"
	"recipeScraper/scrape"
	"strings"
)

func ScrapeRecipe(recipeUrl string) []objects.Recipe {
	u, err := url.Parse(recipeUrl)
	if err != nil {
		log.Fatal(err)
	}

	recipes := make([]objects.Recipe, 0)
	hostname_split := strings.Split(u.Hostname(), ".")
	hostname := hostname_split[len(hostname_split)-2] + "." + hostname_split[len(hostname_split)-1]

	switch hostname {
	case "fitmencook.com":
		recipes = scrape.Fitmencook(recipeUrl)
	case "skinnytaste.com":
		recipes = scrape.Skinnytaste(recipeUrl)
	}

	return recipes
}
