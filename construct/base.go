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
	offset := (pageno - 1) * pagesize
	if pagesize <= 0 || offset < 0 {
		return // 页码或页长错误
	}
	err = query.Count(&total).Limit(pagesize).Offset(offset).Find(out).Error
	if err != nil && gorm.IsRecordNotFoundError(err) {
		err = nil // 忽略没有数据的错误
	}
	return
}
