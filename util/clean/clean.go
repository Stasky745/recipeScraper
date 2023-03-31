package clean

import (
	"regexp"
)

func Strings(s string) string {
	reg := regexp.MustCompile(`\W*([a-zA-Z1-9(].*[a-zA-Z1-9)])\W*`)
	return reg.ReplaceAllString(s, "${1}")
}
