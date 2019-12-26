package models

import (
	"github.com/jinzhu/gorm"
)

var ModelInsts []interface{} // 所有Model实例

// 自动建表，如果缺少表或字段会加上
func MigrateTables(drv string, db *gorm.DB) *gorm.DB {
	return db.AutoMigrate(ModelInsts...) // 创建缺少的表和字段
}