package models

import (
	base "github.com/azhai/gozzo-db/construct"
)

// 查询符合条件的所有行
func (m Access) FindAll(filters ...base.FilterFunc) (objs []*Access, err error) {
	err = db.Model(m).Scopes(filters...).Find(&objs).Error
	err = IgnoreNotFoundError(err)
	return
}

// 查询符合条件的第一行
func (m Access) GetOne(filters ...base.FilterFunc) (obj *Access, err error) {
	obj = new(Access)
	err = db.Model(m).Scopes(filters...).Take(&obj).Error
	err = IgnoreNotFoundError(err)
	return
}
