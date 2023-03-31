package parse

import (
	"log"
	"recipeScraper/objects"
	"regexp"
	"strconv"
	"strings"
)

func Ingredient(s string) objects.Ingredient {
	reg := regexp.MustCompile(`(?i)[^a-zA-Z0-9]*(?:(?:\D*(?P<quantity>\d+(?:\.\d+)?(?:\/\d+)?(?:\s+[½¼]|\s+\d+\/\d+)?)(?:\s?(?:[-–]|or)\s?(?P<extra>\d+)?)?(?:\s?(?:(?P<can>\d+)[-\s])?\s?(?P<unit>\w*\.?)\s+)?(?:(?:can|bottle|of|a)\s+)?)|(?:(?P<pinch>\w+)\s+of))?(?P<ingredient>[^\(]*)(?:\((?P<notes>.*)\))?`)
	// [^a-zA-Z0-9]*(?:\D*(?P<quantity>\d+(?:\.\d+)?(?:\/\d+)?(?:\s+[½¼]|\s+\d+\/\d+)?)(?:\s?(?:[-–]|or)\s?(?P<extra>\d+)?)?(?:\s?(?:(?P<can>\d+)[-\s])?\s?(?P<unit>\w*\.?)\s+)?(?:(?:can|bottle|of|a)\s+)?)?(?P<ingredient>[^\(]*)(?:\((?P<notes>.*)\))?
	// [^a-zA-Z0-9]*(?:\D*(?P<quantity>\d+(?:\.\d+)?(?:\/\d+)?(?:\s+[½¼]|\s+\d+\/\d+)?)(?:\s?(?:[-–]|or)\s?(?P<extra>\d+)?)?)?(?:\s?(?:(?P<can>\d+)[-\s])?\s?(?P<unit>\w*\.?)\s+)?(?:(?:can|bottle|of|a)\s+)?(?P<ingredient>[^\(]*)(?:\((?P<notes>.*)\))?
	match := reg.FindStringSubmatch(s)

	result := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	quantStr := result["quantity"]
	quantStr = strings.ReplaceAll(quantStr, "1/2", ".5")
	quantStr = strings.ReplaceAll(quantStr, ",", ".")
	quantStr = strings.ReplaceAll(quantStr, " ", "")
	quantity := 0.0
	if strings.Contains(quantStr, "¼") {
		quantStr = strings.ReplaceAll(quantStr, "¼", "")
		quantity += 0.25
	} else if strings.Contains(quantStr, "½") {
		quantStr = strings.ReplaceAll(quantStr, "½", "")
		quantity += 0.5
	}
	newQuantity, err := strconv.ParseFloat(quantStr, 64)
	if err != nil {
		newQuantity = 1
	}

	//In case they mention something like a can. Eg: 1 5-ounce can of X
	can, err := strconv.ParseFloat(result["can"], 64)
	if err != nil {
		can = 1
	}
	newQuantity *= can

	quantity += newQuantity

	// this is in case a sentence starts like:
	//    "pinch of X"
	//	  "inch of Y"
	unit := result["unit"]
	if unit == "" && result["pinch"] != "" {
		unit = result["pinch"]
	}

	notes := result["notes"]

	// if the text is something like: 1 or 2 (or with a slash), add the second amount to the notes
	ext, err := strconv.ParseFloat(result["extra"], 64)
	if err == nil {
		if ext > quantity {
			notes = "(or " + strconv.Itoa(int(ext-quantity)) + "more)"
		} else {
			notes = "(or " + strconv.Itoa(int(quantity-ext)) + "less)"
		}
	}

	newIngredient := objects.Ingredient{
		Name:        strings.TrimSpace(result["ingredient"]),
		Quantity:    quantity,
		Unit:        unit,
		Ingredients: nil,
		Notes:       notes,
	}

	return newIngredient
}

func NutritionValue(s string) objects.NutritionValue {
	reg := regexp.MustCompile(`(?P<value>\d+)(?P<unit>\D*)`)
	match := reg.FindStringSubmatch(s)

	result := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	val, err := strconv.Atoi(result["value"])
	if err != nil {
		log.Fatalln("Can't convert value from STR to INT")
	}

	return objects.NutritionValue{
		Unit:  result["unit"],
		Value: val,
	}
}
