package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
)

const DIR_MODE = 0777

// 检查错误
func CheckError(err error) bool {
	if err != nil {
		errInfo := errors.Wrap(err, 1).ErrorStack()
		fmt.Errorf(errInfo)
		return false
	}
	return true
}

// 取得文件的绝对路径
func GetAbsFile(fname string) string {
	if filepath.IsAbs(fname) == false {
		// 相对于程序运行目录
		origDir := filepath.Dir(os.Args[0])
		dir, err := filepath.Abs(origDir)
		if err != nil {
			return ""
		}
		dir = strings.Replace(dir, "\\", "/", -1)
		fname = filepath.Join(dir, fname)
	}
	return fname
}

// detect if file exists
// -1, false 不合法的路径
// 0, false 路径不存在
// -1, true 存在文件夹
// >=0, true 文件并存在
func FileSize(path string) (int64, bool) {
	if path == "" {
		return -1, false
	}
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return 0, false
	}
	var size = int64(-1)
	if info.IsDir() == false {
		size = info.Size()
	}
	return size, true
}

func MkdirForFile(path string) int64 {
	size, exists := FileSize(path)
	if size < 0 {
		return size
	}
	if !exists {
		dir := filepath.Dir(path)
		_ = os.MkdirAll(dir, DIR_MODE)
	}
	return size
}
