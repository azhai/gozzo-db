package construct

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// 字符串比较方式
const (
	CMP_STRING_OMIT             = iota // 不比较
	CMP_STRING_CONTAINS                // 包含
	CMP_STRING_STARTSWITH              // 打头
	CMP_STRING_ENDSWITH                // 结尾
	CMP_STRING_IGNORE_SPACES           // 忽略空格
	CMP_STRING_CASE_INSENSITIVE        // 不分大小写
	CMP_STRING_EQUAL                   // 相等
)

// 删除所有空白，包括中间的
func RemoveSpaces(s string) string {
	subs := map[string]string{
		" ":"", "\n":"", "\r":"", "\t":"", "\v":"", "\f":"",
	}
	return ReplaceWith(s, subs)
}


// 一一对应进行替换，次序不定（因为map的关系）
func ReplaceWith(s string, subs map[string]string) string {
	if s == "" {
		return ""
	}
	var marks []string
	for key, value := range subs {
		marks = append(marks, key, value)
	}
	replacer := strings.NewReplacer(marks...)
	return replacer.Replace(s)
}


// 比较是否相符
func StringMatch(a, b string, cmp int) bool {
	switch cmp {
	case CMP_STRING_OMIT:
		return true
	case CMP_STRING_CONTAINS:
		return strings.Contains(a, b)
	case CMP_STRING_STARTSWITH:
		return strings.HasPrefix(a, b)
	case CMP_STRING_ENDSWITH:
		return strings.HasSuffix(a, b)
	case CMP_STRING_IGNORE_SPACES:
		a, b = RemoveSpaces(a), RemoveSpaces(b)
		return strings.EqualFold(a, b)
	case CMP_STRING_CASE_INSENSITIVE:
		return strings.EqualFold(a, b)
	default: // 包括 CMP_STRING_EQUAL
		return strings.Compare(a, b) == 0
	}
}

// 是否在字符串列表中
func InStringList(x string, lst []string, cmp int) bool {
	size := len(lst)
	if size == 0 {
		return false
	}
	if !sort.StringsAreSorted(lst) {
		sort.Strings(lst)
	}
	i := sort.Search(size, func(i int) bool { return lst[i] >= x })
	return i < size && StringMatch(x, lst[i], cmp)
}


// 是否在字符串列表中，比较方式是有任何一个开头符合
func StartStringList(x string, lst []string) bool {
	return InStringList(x, lst, CMP_STRING_STARTSWITH)
}

func IsSubsetList(lst1, lst2 []string) bool {
	for _, x := range lst1 {
		if !InStringList(x, lst2, CMP_STRING_EQUAL) {
			return false
		}
	}
	return true
}

// 遍历目录下的文件
func FindFiles(dir, ext string) (map[string]os.FileInfo, error) {
	var result = make(map[string]os.FileInfo)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return result, err
	}
	for _, file := range files {
		fname := file.Name()
		if ext != "" && !strings.HasSuffix(fname, ext) {
			continue
		}
		fname = filepath.Join(dir, fname)
		result[fname] = file
	}
	return result, nil
}


