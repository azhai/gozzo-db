package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/azhai/gozzo-db/prepare"
	"github.com/azhai/gozzo-db/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	diaName  string // 数据库连接名
	fileName string // 配置文件名
	verbose  bool   // 详细输出
)

func init() {
	flag.StringVar(&diaName, "d", "default", "数据库连接名")
	flag.StringVar(&fileName, "f", "settings.toml", "配置文件名")
	flag.BoolVar(&verbose, "v", false, "输出详细信息")
	flag.Parse()
}

func main() {
	conf := prepare.GetConfig(fileName)
	if verbose {
		fmt.Printf("%+v\n\n", conf)
	}
	drv, dsn := conf.GetDSN(diaName)
	db, err := gorm.Open(drv, dsn)
	if !utils.CheckError(err) {
		return
	}
	if verbose {
		db.LogMode(true).SetLogger(log.New(os.Stdout, "\r\n", 0))
	}
	names, _ := prepare.CreateModels(db, conf.Application)
	if drv == "sqlite3" {
		drv = "sqlite" // Sqlite的import包名和Open()驱动名不一样
	}
	_ = prepare.GenInitFile(names, conf.Application, fileName, diaName, drv)
}
