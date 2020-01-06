package main

import (
	"flag"
	"fmt"

	"github.com/azhai/gozzo-db/cmd"
	"github.com/azhai/gozzo-db/cmd/models"
	"github.com/azhai/gozzo-db/prepare"
)

var (
	onlyTableComment    bool   // 只修改表注释
	targetDir, baseName string // 临时目录和基础Model名
	verbose             bool   // 详细输出
)

// NOTE: 编译依赖于gen2model生成的models
// 从代码中同步到数据表中，包括缺少的字段、索引和改动的注释
func main() {
	flag.BoolVar(&onlyTableComment, "tc", false, "只修改表注释")
	flag.StringVar(&targetDir, "td", "cmd/tmp/", "临时目录")
	flag.StringVar(&baseName, "bn", "BaseModel", "基础Model名")
	flag.BoolVar(&verbose, "vv", false, "输出详细信息")
	conf := cmd.Initialize(nil)
	db, err := cmd.ConnectDatabase(conf, conf.ConnName)
	if err != nil || db == nil {
		panic(err)
	}

	// 创建表结构
	drv := conf.GetDriverName("default")
	db = models.MigrateTables(drv, db)
	// 更新字段或表注释
	opts := prepare.Options{
		SourceDir:        conf.Application.OutputDir,
		TargetDir:        targetDir,
		BaseName:         baseName,
		OnlyTableComment: onlyTableComment,
	}
	opts.TablePrefix = conf.GetTablePrefix(conf.ConnName)
	err = prepare.AmendComments(db.DB(), opts, verbose)
	if err != nil {
		fmt.Println(err)
		return
	}
}
