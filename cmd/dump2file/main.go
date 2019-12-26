package main

import (
	"flag"
	"fmt"

	"github.com/azhai/gozzo-db/cmd"
	"github.com/azhai/gozzo-db/cmd/models"
	"github.com/azhai/gozzo-db/export"
)

var (
	outFile, srcFile    string // 数据源文件与输出文件
	outLimit            int    // 每张表输出前几行
	verbose             bool   // 详细输出
)

// NOTE: 编译依赖于gen2model生成的models
// 在数据表和TOML文件之间导入导出数据
func main() {
	flag.StringVar(&outFile, "of", "", "默认操作，输出文件，为空时使用配置文件中的data_file")
	flag.StringVar(&srcFile, "sf", "", "数据源文件")
	flag.IntVar(&outLimit, "ol", -1, "每张表输出前几行")
	flag.BoolVar(&verbose, "vv", false, "输出详细信息")
	flag.Parse()

	conf, db := cmd.Initialize()
	if outFile == "" {
		outFile = conf.Application.DataFile
	}
	if srcFile != "" { // 导入数据
		_, err := export.LoadFileData(db, srcFile, models.ModelInsts, verbose)
		if err != nil {
			fmt.Println(err)
		}
	}
	if outFile != "" { // 导出数据
		ep := export.NewExportor()
		for _, m := range models.ModelInsts {
			ep.AddQueryResult(m, db.Limit(outLimit))
		}
		if err := ep.WriteTo(outFile); err != nil {
			fmt.Println(err)
		}
	}
}
