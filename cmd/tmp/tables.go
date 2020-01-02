package tmp

import (
	"time"

	base "github.com/azhai/gozzo-db/construct"
)

type BaseModel = base.Model

// 权限控制
type Access struct {
	BaseModel
	RoleName     string     `json:"role_name" gorm:"size:50;not null;index;default:'';comment:'角色名'"` // 角色名
	ResourceType string     `json:"resource_type" gorm:"size:50;not null;default:'';comment:'资源类型'"`  // 资源类型
	ResourceArgs *string    `json:"resource_args" gorm:"comment:'资源参数'"`                              // 资源参数
	PermCode     uint       `json:"perm_code" gorm:"not null;default:'0';comment:'权限码'"`              // 权限码
	Actions      string     `json:"actions" gorm:"size:50;not null;default:'';comment:'允许的操作'"`       // 允许的操作
	GrantedAt    *time.Time `json:"granted_at" gorm:"comment:'授权时间'"`                                 // 授权时间
	RevokedAt    *time.Time `json:"revoked_at" gorm:"index;comment:'撤销时间'"`                           // 撤销时间
}

// 数据表名为 t_access
func (Access) TableName() string {
	return "t_access"
}

// 数据表备注
func (Access) TableComment() string {
	return "权限控制"
}

// 用户组
type Group struct {
	BaseModel
	GID       string    `json:"gid" gorm:"unique_index;type:char(16);not null;default:'';comment:'唯一ID';column:gid"` // 唯一ID
	Title     string    `json:"title" gorm:"size:50;not null;default:'';comment:'名称'"`                               // 名称
	Remark    *string   `json:"remark" gorm:"type:text;comment:'说明备注'"`                                              // 说明备注
	CreatedAt time.Time `json:"CreatedAt" gorm:"comment:'创建时间'"`                                                     // 创建时间
}

// 数据表名为 t_group
func (Group) TableName() string {
	return "t_group"
}

// 数据表备注
func (Group) TableComment() string {
	return "用户组"
}

// 菜单
type Menu struct {
	BaseModel
	Lft       uint       `json:"lft" gorm:"not null;default:'0';comment:'左边界'"`               // 左边界
	Rgt       uint       `json:"rgt" gorm:"not null;index;default:'0';comment:'右边界'"`         // 右边界
	Depth     uint8      `json:"depth" gorm:"not null;index;default:'1';comment:'高度'"`        // 高度
	Path      string     `json:"path" gorm:"size:100;not null;index;default:'';comment:'路径'"` // 路径
	Title     string     `json:"title" gorm:"size:50;not null;default:'';comment:'名称'"`       // 名称
	Icon      *string    `json:"icon" gorm:"size:30;comment:'图标'"`                            // 图标
	Remark    *string    `json:"remark" gorm:"type:text;comment:'说明备注'"`                      // 说明备注
	CreatedAt time.Time  `json:"CreatedAt" gorm:"comment:'创建时间'"`                             // 创建时间
	UpdatedAt time.Time  `json:"UpdatedAt" gorm:"comment:'更新时间'"`                             // 更新时间
	DeletedAt *time.Time `json:"DeletedAt" gorm:"index;comment:'删除时间'"`                       // 删除时间
}

// 数据表名为 t_menu
func (Menu) TableName() string {
	return "t_menu"
}

// 数据表备注
func (Menu) TableComment() string {
	return "菜单"
}

// 角色
type Role struct {
	BaseModel
	Name      string     `json:"name" gorm:"unique_index;size:50;not null;default:'';comment:'名称'"` // 名称
	Remark    *string    `json:"remark" gorm:"type:text;comment:'说明备注'"`                            // 说明备注
	CreatedAt time.Time  `json:"CreatedAt" gorm:"comment:'创建时间'"`                                   // 创建时间
	UpdatedAt time.Time  `json:"UpdatedAt" gorm:"comment:'更新时间'"`                                   // 更新时间
	DeletedAt *time.Time `json:"DeletedAt" gorm:"index;comment:'删除时间'"`                             // 删除时间
}

// 数据表名为 t_role
func (Role) TableName() string {
	return "t_role"
}

// 数据表备注
func (Role) TableComment() string {
	return "角色"
}

// 用户
type User struct {
	BaseModel
	UID          string     `json:"uid" gorm:"unique_index;type:char(16);not null;default:'';comment:'唯一ID'"` // 唯一ID
	Username     string     `json:"username" gorm:"size:30;not null;index;default:'';comment:'用户名'"`          // 用户名
	Password     string     `json:"-" gorm:"size:60;not null;default:'';comment:'密码'"`                        // 密码
	Realname     *string    `json:"realname" gorm:"size:20;comment:'昵称/称呼'"`                                  // 昵称/称呼
	Mobile       *string    `json:"mobile" gorm:"size:20;index;comment:'手机号码'"`                               // 手机号码
	Email        *string    `json:"email" gorm:"size:50;comment:'电子邮箱'"`                                      // 电子邮箱
	PrinGid      string     `json:"prin_gid" gorm:"type:char(16);not null;default:'';comment:'主用户组'"`         // 主用户组
	ViceGid      *string    `json:"vice_gid" gorm:"type:char(16);comment:'次用户组'"`                             // 次用户组
	Avatar       *string    `json:"avatar" gorm:"size:100;comment:'头像'"`                                      // 头像
	Introduction *string    `json:"introduction" gorm:"size:500;comment:'介绍说明'"`                              // 介绍说明
	CreatedAt    time.Time  `json:"CreatedAt" gorm:"comment:'创建时间'"`                                          // 创建时间
	UpdatedAt    time.Time  `json:"UpdatedAt" gorm:"comment:'更新时间'"`                                          // 更新时间
	DeletedAt    *time.Time `json:"DeletedAt" gorm:"index;comment:'删除时间'"`                                    // 删除时间
}

// 数据表名为 t_user
func (User) TableName() string {
	return "t_user"
}

// 数据表备注
func (User) TableComment() string {
	return "用户"
}

// 用户角色
type UserRole struct {
	BaseModel
	UserUID  string `json:"user_uid" gorm:"type:char(16);not null;index;default:'';comment:'用户ID'"` // 用户ID
	RoleName string `json:"role_name" gorm:"size:50;not null;index;default:'';comment:'角色名'"`       // 角色名
}

// 数据表名为 t_user_role
func (UserRole) TableName() string {
	return "t_user_role"
}

// 数据表备注
func (UserRole) TableComment() string {
	return "用户角色"
}
