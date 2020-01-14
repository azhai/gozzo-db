package models

import (
	base "github.com/azhai/gozzo-db/construct"
)

// 查询符合条件的所有行
func (m User) FindAll(filters ...base.FilterFunc) (objs []*User, err error) {
	err = db.Model(m).Scopes(filters...).Find(&objs).Error
	err = IgnoreNotFoundError(err)
	return
}

// 查询符合条件的第一行
func (m User) GetOne(filters ...base.FilterFunc) (obj *User, err error) {
	obj = new(User)
	err = db.Model(m).Scopes(filters...).Take(&obj).Error
	err = IgnoreNotFoundError(err)
	return
}
