package prepare

import (
	"fmt"
	"io"
	"log"

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

// gorm 之前版本（v1.9.11）日志没有在各部分中间加空格
type Logger struct {
	Space string
	*log.Logger
}

func NewLogger(out io.Writer, space string) *Logger {
	logger := log.New(out, "\r\n", log.LstdFlags)
	return &Logger{space, logger}
}

func (l *Logger) Print(v ...interface{}) {
	var vv []interface{}
	for _, x := range v {
		vv = append(vv, x, l.Space)
	}
	l.Output(2, fmt.Sprint(vv...))
}
