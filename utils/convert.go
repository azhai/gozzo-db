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

func TrimTail(lines string) string {
	return strings.Join(strings.Fields(lines), " ")
}

func WrapWith(name, left, right string) string {
	if name == "" {
		return ""
	}
	return fmt.Sprintf("%s%s%s", left, name, right)
}
