package construct

import (
	"github.com/jinzhu/gorm"
)

type Model struct {
	ID uint `json:"id" gorm:"primary_key;not null;auto_increment"`
}

/**
 * 过滤查询
 * 使用方法 query = query.Scopes(filters ...FilterFunc)
 */
type FilterFunc = func(query *gorm.DB) *gorm.DB

/**
 * 翻页查询，out参数需要传引用
 * 使用方法 total, err := Paginate(query, &rows, pageno, pagesize)
 */
func Paginate(query *gorm.DB, out interface{}, pageno, pagesize int) (total int, err error) {
	query = query.Count(&total)
	offset, limit := -1, -1 // 这也是gorm的初始值
	// 参数校正
	if pagesize >= 0 {
		limit = pagesize
		if page >= 0 {
			offset = (pageno - 1) * pagesize
		} else if page < 0 {
			offset = total + pageno * pagesize
		}
	}
	if offset < 0 {
		offset = -1
	}
	err = query.Limit(limit).Offset(offset).Find(out).Error
	if err != nil && gorm.IsRecordNotFoundError(err) {
		err = nil // 忽略没有数据的错误
	}
	return
}
