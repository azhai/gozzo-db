package main

import (
	"flag"
	"fmt"

	"github.com/azhai/gozzo-db/cmd"
	"github.com/azhai/gozzo-db/cmd/models"
	"github.com/azhai/gozzo-db/export"
	"github.com/azhai/gozzo-db/prepare"
)

var (
	srcFile          string // 数据源文件
	outFile          string // 输出文件
	outLimit         int    // 每张表输出前几行
	onlyTableComment bool   // 只修改表注释
	verbose          bool   // 详细输出
)

// NOTE: 编译依赖于table2file生成的models
func main() {
	flag.StringVar(&srcFile, "sf", "", "数据源文件")
	flag.StringVar(&outFile, "of", "", "输出文件")
	flag.IntVar(&outLimit, "ol", -1, "每张表输出前几行")
	flag.BoolVar(&onlyTableComment, "tc", false, "只修改表注释")
	flag.BoolVar(&verbose, "vv", false, "输出详细信息")
	flag.Parse()

	conf, db := cmd.Initialize()
	if outFile == "" {
		outFile = conf.Application.DataFile
	}

	// 创建表结构
	drv := conf.GetDriverName("default")
	db = models.MigrateTables(drv, db)
	// 更新字段或表注释
	pre := conf.GetTablePrefix(conf.ConnName)
	err := prepare.AmendComments(db.DB(), pre, onlyTableComment, verbose)
	if err != nil {
		fmt.Println(err)
		return
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
