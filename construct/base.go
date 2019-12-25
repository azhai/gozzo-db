package construct

import (
	"github.com/jinzhu/gorm"
)

/**
 * 过滤查询
 * 使用方法 query = query.Scopes(filters ...FilterFunc)
 */
type FilterFunc = func(query *gorm.DB) *gorm.DB

/**
 * 数据表名
 */
type ITableName interface {
	TableName() string
}

/**
 * 数据表注释
 */
type ITableComment interface {
	TableComment() string
}

/**
 * 带自增主键的基础Model
 */
type Model struct {
	ID uint `json:"id" toml:"-" gorm:"primary_key;not null;auto_increment"`
}

func (Model) TableComment() string {
	return ""
}

// 忽略表中无数据的错误
func IgnoreNotFoundError(err error) error {
	if err == nil || gorm.IsRecordNotFoundError(err) {
		return nil
	}
	return err
}

/**
 * 翻页查询，out参数需要传引用
 * 使用方法 total, err := Paginate(query, &rows, pageno, pagesize)
 */
func Paginate(query *gorm.DB, out interface{}, pageno, pagesize int) (total int, err error) {
	err = query.Count(&total).Error
	if err != nil || total <= 0 {
		return
	}
	offset, limit := -1, -1 // 这也是gorm的初始值
	// 参数校正
	if pagesize >= 0 {
		limit = pagesize
		if pageno >= 0 {
			offset = (pageno - 1) * pagesize
		} else if pageno < 0 {
			offset = total + pageno*pagesize
		}
	}
	if offset < 0 {
		offset = -1
	}
	err = query.Limit(limit).Offset(offset).Find(out).Error
	err = IgnoreNotFoundError(err)
	return
}
