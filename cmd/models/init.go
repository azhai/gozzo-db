package models

import (
	"log"
	"os"

	"github.com/azhai/gozzo-db/cache"
	base "github.com/azhai/gozzo-db/construct"
	"github.com/azhai/gozzo-db/export"
	"github.com/azhai/gozzo-db/prepare"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db         *gorm.DB         // 数据库对象
	ModelInsts = []interface{}{ // 所有Model实例
		&Access{}, &Group{}, &Menu{},
		&Role{}, &User{}, &UserRole{},
	}
)

type BaseModel = base.Model

// 忽略表中无数据的错误
func IgnoreNotFoundError(err error) error {
	return base.IgnoreNotFoundError(err)
}

// 获取当前db
func Query() *gorm.DB {
	return db
}

// 查询某张数据表
func QueryTable(name string) *gorm.DB {
	return db.Table(name)
}

// 连接数据库
func init() {
	conf, err := prepare.GetConfig("settings.toml")
	if err != nil {
		panic(err)
	}
	if c, ok := conf.Connections["cache"]; ok && c.Driver == "redis" {
		rds := cache.ConnectRedisPool(c.ConnParams)
		cache.SetRedisPool(rds)
	}
	db, err = gorm.Open(conf.GetDSN("default"))
	if err != nil {
		panic(err)
	}

	// 初始化数据库
	if conf.Application.Debug {
		db = db.Debug().LogMode(true)
		db.SetLogger(log.New(os.Stdout, "\r\n", 0))
	}
	drv := conf.GetDriverName("default")
	if drv == "mysql" {
		db.Set("gorm:table_options", "ENGINE=InnoDB")
	}
	db = MigrateTables(drv, db)
}

// 自动建表，如果缺少表或字段会加上
func MigrateTables(drv string, db *gorm.DB) *gorm.DB {
	db.SingularTable(true)
	db = db.AutoMigrate(ModelInsts...) // 创建缺少的表和字段
	if drv == "mysql" {                // 更新MySQL表注释
		db = export.AlterTableComments(db, ModelInsts...)
	}
	return db
}
