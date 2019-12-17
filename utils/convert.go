package utils

import (
	"fmt"
	"strings"

	"github.com/codemodus/kace"
	"github.com/jinzhu/inflection"
)

func ToCamel(name string) string {
	return kace.Pascal(name)
}

func ToSnake(name string) string {
	return kace.Snake(name)
}

func ToPlural(name string) string {
	return inflection.Plural(name)
}

func ToSingular(name string) string {
	return inflection.Singular(name)
}

func ReduceSpaces(lines string) string {
	return strings.Join(strings.Fields(lines), " ")
}

func WrapWith(s, left, right string) string {
	if s == "" {
		return ""
	}
	return fmt.Sprintf("%s%s%s", left, s, right)
}

func ReplaceQuotes(s string) string {
	if s == "" {
		return ""
	}
	replacer := strings.NewReplacer("[", "`", "]", "`")
	return replacer.Replace(s)
}
