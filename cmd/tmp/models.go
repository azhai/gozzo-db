package tmp

import (
	"time"

	base "github.com/azhai/gozzo-db/construct"
	"github.com/jinzhu/gorm"
)


// 用户组
type Group struct {
	BaseModel
	Gid       string    `json:"gid" gorm:"size:16;not null;comment:'唯一ID'"`               // 唯一ID
	Title     string    `json:"title" gorm:"size:50;not null;comment:'名称'"`               // 名称
	Remark    *string   `json:"remark" gorm:"comment:'说明备注'"`                             // 说明备注
	GID       string    `json:"g_id" gorm:"unique_index;size:16;not null;comment:'唯一ID'"` // 唯一ID
	CreatedAt time.Time `json:"created_at" gorm:"comment:'创建时间'"`                         // 创建时间
}

// 菜单
type Menu struct {
	BaseModel
	Title     string     `json:"title" gorm:"size:50;not null;comment:'名称'"`    // 名称
	Remark    *string    `json:"remark" gorm:"comment:'说明备注'"`                  // 说明备注
	LeftEdge  uint       `json:"left_edge" gorm:"not null;index;comment:'左边界'"` // 左边界
	RightEdge uint       `json:"right_edge" gorm:"not null;comment:'右边界'"`      // 右边界
	SeqNo     uint8      `json:"seq_no" gorm:"not null;comment:'次序'"`           // 次序
	CreatedAt time.Time  `json:"created_at" gorm:"comment:'创建时间'"`              // 创建时间
	UpdatedAt time.Time  `json:"updated_at" gorm:"comment:'更新时间'"`              // 更新时间
	DeletedAt *time.Time `json:"deleted_at" gorm:"index;comment:'删除时间'"`        // 删除时间
}

// 拥有权限
type Permission struct {
	BaseModel
	RoleID       uint    `json:"role_id" gorm:"not null;index;comment:'角色ID'"`         // 角色ID
	ResourceType string  `json:"resource_type" gorm:"size:50;not null;comment:'资源类型'"` // 资源类型
	ResourceArgs *string `json:"resource_args" gorm:"comment:'资源参数'"`                  // 资源参数
	Action       string  `json:"action" gorm:"size:20;not null;comment:'允许的操作'"`       // 允许的操作
}

// 角色
type Role struct {
	BaseModel
	Title     string     `json:"title" gorm:"size:50;not null;comment:'名称'"` // 名称
	Remark    *string    `json:"remark" gorm:"comment:'说明备注'"`               // 说明备注
	CreatedAt time.Time  `json:"created_at" gorm:"comment:'创建时间'"`           // 创建时间
	UpdatedAt time.Time  `json:"updated_at" gorm:"comment:'更新时间'"`           // 更新时间
	DeletedAt *time.Time `json:"deleted_at" gorm:"index;comment:'删除时间'"`     // 删除时间
}

// 用户
type User struct {
	BaseModel
	UID       string     `json:"uid" gorm:"unique_index;size:16;not null;comment:'唯一ID'"` // 唯一ID
	Username  string     `json:"username" gorm:"size:30;not null;index;comment:'用户名'"`    // 用户名
	Password  string     `json:"-" gorm:"size:60;not null;comment:'密码'"`                  // 密码
	Realname  *string    `json:"realname" gorm:"size:20;comment:'昵称/称呼'"`                 // 昵称/称呼
	Mobile    *string    `json:"mobile" gorm:"size:20;index;comment:'手机号码'"`              // 手机号码
	Email     *string    `json:"email" gorm:"size:50;comment:'电子邮箱'"`                     // 电子邮箱
	PrinGid   string     `json:"prin_gid" gorm:"size:16;not null;comment:'主用户组'"`         // 主用户组
	ViceGid   *string    `json:"vice_gid" gorm:"size:16;comment:'次用户组'"`                  // 次用户组
	CreatedAt time.Time  `json:"created_at" gorm:"comment:'创建时间'"`                        // 创建时间
	UpdatedAt time.Time  `json:"updated_at" gorm:"comment:'更新时间'"`                        // 更新时间
	DeletedAt *time.Time `json:"deleted_at" gorm:"index;comment:'删除时间'"`                  // 删除时间
}

// 用户角色
type UserRole struct {
	BaseModel
	UserID uint `json:"user_id" gorm:"not null;index"`
	RoleID uint `json:"role_id" gorm:"not null"`
}
