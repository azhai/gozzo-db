package models

import (
	"github.com/azhai/gozzo-db/export"
	"github.com/jinzhu/gorm"
)

var ModelInsts []interface{} // 所有Model实例

// 自动建表，如果缺少表或字段会加上
func MigrateTables(drv string, db *gorm.DB) *gorm.DB {
	db = db.AutoMigrate(ModelInsts...) // 创建缺少的表和字段
	if drv == "mysql" { // 更新MySQL表注释
		db = export.AlterTableComments(db, ModelInsts...)
	}
	return db
}

// 写入必须的初始化数据
func FillRequiredData(drv string, db *gorm.DB) *gorm.DB {
	return db
}