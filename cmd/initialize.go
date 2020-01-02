package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/azhai/gozzo-db/prepare"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	connName string // 数据库连接名
	fileName string // 配置文件名
	verbose  bool   // 详细输出
)

func ConnectDatabase(conf *prepare.Config, name string) (*gorm.DB, error) {
	// 连接数据库生成models
	db, err := gorm.Open(conf.GetDSN(name))
	if err != nil {
		return nil, err
	}
	if verbose {
		db = db.Debug().LogMode(true)
		db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	}
	return db, nil
}

// 初始化，解析配置和连接数据库
func Initialize(parse func() error) (*prepare.Config, *gorm.DB) {
	flag.StringVar(&connName, "d", "default", "数据库连接名")
	flag.StringVar(&fileName, "f", "settings.toml", "配置文件名")
	flag.BoolVar(&verbose, "v", false, "输出详细信息")
	flag.Parse()
	if parse != nil {
		if err := parse(); err != nil {
			fmt.Println(err.Error())
			return nil, nil
		}
	}

	// 解析配置文件
	conf, err := prepare.GetConfig(fileName)
	if verbose {
		fmt.Printf("%s:\n%+v\n\n", fileName, conf)
	}
	if err != nil || conf == nil {
		panic(err)
	}
	conf.ConnName = connName
	var db *gorm.DB
	db, err = ConnectDatabase(conf, conf.ConnName)
	if err != nil || db == nil {
		panic(err)
	}
	return conf, db
}
