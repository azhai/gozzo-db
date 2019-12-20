package main

import (
	"flag"
	"fmt"

	"github.com/azhai/gozzo-db/cmd"
	"github.com/azhai/gozzo-db/export"
	"github.com/azhai/gozzo-db/models"
)

var (
	srcFile string // 数据源文件
	outFile string // 输出文件
	verbose bool   // 详细输出
)

func main() {
	flag.StringVar(&srcFile, "src", "cmd/tmp/data.toml", "数据源文件")
	flag.StringVar(&outFile, "o", "", "输出文件")
	flag.BoolVar(&verbose, "vv", false, "输出详细信息")
	flag.Parse()

	var err error
	_, db := cmd.Initialize()
	if outFile == "" { // 导入数据
		_, err = export.LoadFileData(db, srcFile, models.ModelInsts, verbose)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	// 导出数据
	ep := export.NewExportor()
	for _, m := range models.ModelInsts {
		var table string
		if table = export.GetTableName(m, ""); table == "" {
			continue
		}
		var objs []interface{}
		err = db.Model(m).Table(table).Find(&objs).Error
		if size := len(objs); size > 0 {
			fmt.Sprintf("%s %d\n", table, size)
			/*for i := 0; i < size; i ++ {
				ep.Data[table] = append(ep.Data[table], objs[i])
			}*/
			// ep.Data[table] = objs
		}
	}
	if err = ep.WriteTo(outFile); err != nil {
		fmt.Println(err)
	}
}
