package main

import (
	"flag"
	"fmt"

	"github.com/azhai/gozzo-db/cmd"
	"github.com/azhai/gozzo-db/prepare"
)

var (
	onlyTable bool // 只修改表注释
	verbose   bool // 详细输出
)

func main() {
	flag.BoolVar(&onlyTable, "ot", false, "只修改表注释")
	flag.BoolVar(&verbose, "vv", false, "输出详细信息")
	flag.Parse()

	conf, db := cmd.Initialize()
	tablePrefix := conf.GetTablePrefix(conf.ConnName)
	err := prepare.AmendComments(db.DB(), tablePrefix, onlyTable, verbose)
	if err != nil {
		fmt.Println(err)
	}
}
