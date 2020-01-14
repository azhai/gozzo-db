package models

import (
	base "github.com/azhai/gozzo-db/construct"
)

// 查询符合条件的所有行
func (m UserRole) FindAll(filters ...base.FilterFunc) (objs []*UserRole, err error) {
	err = db.Model(m).Scopes(filters...).Find(&objs).Error
	err = IgnoreNotFoundError(err)
	return
}

// 查询符合条件的第一行
func (m UserRole) GetOne(filters ...base.FilterFunc) (obj *UserRole, err error) {
	obj = new(UserRole)
	err = db.Model(m).Scopes(filters...).Take(&obj).Error
	err = IgnoreNotFoundError(err)
	return
}
