package models

import (
	base "github.com/azhai/gozzo-db/construct"
)

// 查询符合条件的所有行
func (m Role) FindAll(filters ...base.FilterFunc) (objs []*Role, err error) {
	err = db.Model(m).Scopes(filters...).Find(&objs).Error
	err = IgnoreNotFoundError(err)
	return
}

// 查询符合条件的第一行
func (m Role) GetOne(filters ...base.FilterFunc) (obj *Role, err error) {
	obj = new(Role)
	err = db.Model(m).Scopes(filters...).Take(&obj).Error
	err = IgnoreNotFoundError(err)
	return
}
