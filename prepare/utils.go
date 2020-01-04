package prepare

import (
	"fmt"
	"github.com/azhai/gozzo-utils/filesystem"
	"github.com/codemodus/kace"
	"github.com/go-errors/errors"
	"github.com/jinzhu/inflection"
	"os"
	"path/filepath"
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

// 为文件路径创建目录
func MkdirForFile(path string) int64 {
	size, exists := filesystem.FileSize(path)
	if size < 0 {
		return size
	}
	if !exists {
		dir := filepath.Dir(path)
		_ = os.MkdirAll(dir, DIR_MODE)
	}
	return size
}
