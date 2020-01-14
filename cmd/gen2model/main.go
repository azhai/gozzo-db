package main

import (
	"flag"
	"fmt"

	"github.com/azhai/gozzo-db/cmd"
	"github.com/azhai/gozzo-db/prepare"
	"github.com/azhai/gozzo-utils/filesystem"
)

/**
 * 生成文件不同结构，可选 0-5
 * 0  与 5 类似，但会在每个文件名前面加一个下划线
 * 1  只生成 init.go 文件
 * 2  除了 init.go 文件， table 和 query 都放入 tables.go 中
 * 3  除了 init.go 文件， table 都放入 tables.go 中， query 都放入 queries.go 中
 * 4  除了 init.go 文件， table 都放入 tables.go 中， query 分开放入对应模型名文件中
 * 5  除了 init.go 文件， table 和 query 一起放入对应模型名文件中
 */
var mode uint

// 根据数据表结构生成对应的Model代码
func main() {
	flag.UintVar(&mode, "mode", 0, "生成文件不同结构，可选 0-5 ，具体参考说明")
	conf := cmd.Initialize(func() error {
		if mode >= 6 {
			return fmt.Errorf("没有这种模式，可选 0-5 ，具体参考说明")
		}
		return nil
	})
	db, err := cmd.ConnectDatabase(conf, conf.ConnName)
	if err != nil || db == nil {
		panic(err)
	}

	var names []string
	outDir := conf.Application.OutputDir
	filesystem.MkdirForFile(fmt.Sprintf("%s/init.go", outDir))
	names, err = prepare.CreateModels(conf, db, mode)
	if err != nil {
		fmt.Println(err)
	}
	err = prepare.GenInitFile(conf, names, mode)
	if err != nil {
		fmt.Println(err)
	}
}
