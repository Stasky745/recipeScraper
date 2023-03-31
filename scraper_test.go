package recipeScraper

import (
	"recipeScraper/test"
	"testing"
)

func TestScrapeRecipe(t *testing.T) {
	// fitmencook.com
	test.FitMenCookMultiRecipe(t, ScrapeRecipe("https://fitmencook.com/5-air-fried-chicken-recipes/"))
	test.FitMenCookSingleRecipe(t, ScrapeRecipe("https://fitmencook.com/banana-overnight-oats/"))

	// skinnytaste.com
	test.SkinnyTasteSimpleIngredients(t, ScrapeRecipe("https://www.skinnytaste.com/spicy-canned-salmon-rice-bowl/"))
	test.SkinnyTasteMultiIngredients(t, ScrapeRecipe("https://www.skinnytaste.com/slow-cooker-banh-mi-rice-bowls/"))
}
