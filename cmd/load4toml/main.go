package main

import (
	"flag"
	"fmt"

	"github.com/azhai/gozzo-db/cmd"
	"github.com/azhai/gozzo-db/export"
	"github.com/azhai/gozzo-db/models"
)

var (
	outFile string // 输出文件
	verbose bool   // 详细输出
)

func main() {
	flag.StringVar(&outFile, "of", "cmd/tmp/data.toml", "输出文件")
	flag.BoolVar(&verbose, "vv", false, "输出详细信息")
	flag.Parse()

	_, db := cmd.Initialize()
	_, err := export.LoadFileData(db, outFile, models.ModelInsts, verbose)
	if err != nil {
		fmt.Println(err)
	}
}
