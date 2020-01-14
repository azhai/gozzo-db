package prepare

import (
	"fmt"

	"github.com/codemodus/kace"
	"github.com/go-errors/errors"
	"github.com/jinzhu/inflection"
)

const DIR_MODE = 0777

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

// 检查错误
func CheckError(err error) bool {
	if err != nil {
		errInfo := errors.Wrap(err, 1).ErrorStack()
		fmt.Errorf(errInfo)
		return false
	}
	return true
}
