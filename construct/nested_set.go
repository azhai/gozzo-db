package construct

import (
	"github.com/jinzhu/gorm"
)

// 嵌套集合树
type NestedNode struct {
	Lft   uint  `json:"-" gorm:"not null;index;default:0;comment:'左边界'"` // 左边界
	Rgt   uint  `json:"-" gorm:"not null;index;default:0;comment:'右边界'"` // 右边界
	Depth uint8 `json:"depth" gorm:"not null;default:0;comment:'高度'"`    // 高度
}

// 是否叶子节点
func (n NestedNode) IsLeaf() bool {
	return n.Rgt-n.Lft == 1
}

// 有多少个子孙节点
func (n NestedNode) CountChildren() int {
	return int(n.Rgt-n.Lft-1) / 2
}

// 找出所有祖先节点
func (n NestedNode) AncestorsFilter(Backward bool) FilterFunc {
	return func(query *gorm.DB) *gorm.DB {
		query = query.Where("rgt > ? AND lft < ?", n.Rgt, n.Lft)
		if Backward { // 从子孙往祖先方向排序，即时间倒序
			return query.Order("rgt ASC")
		} else {
			return query.Order("rgt DESC")
		}
	}
}

// 找出所有子孙节点
func (n NestedNode) ChildrenFilter(rank uint8) FilterFunc {
	return func(query *gorm.DB) *gorm.DB {
		if n.Rgt > 0 && n.Lft > 0 { // 当前不是第0层，即具体某分支以下的节点
			query = query.Where("rgt < ? AND lft > ?", n.Rgt, n.Lft)
		}
		if rank > 0 { // 限制层级
			query = query.Where("depth < ?", n.Depth+rank)
		}
		if rank == 1 { // 只有一层，两种排序等价，尽量使用rgt作为索引
			return query.Order("rgt ASC")
		} else {
			return query.Order("lft ASC")
		}
	}
}

// 添加到父节点最末，query一定要使用db.Table(...)
func (n *NestedNode) AddToParent(parent *NestedNode, query *gorm.DB) error {
	var finder = query.Order("rgt DESC")
	if parent == nil {
		n.Depth = 1
	} else {
		n.Depth = parent.Depth + 1
		finder = finder.Where("rgt < ? AND lft > ?", parent.Rgt, parent.Lft)
	}
	finder = finder.Where("depth = ?", n.Depth)
	sibling := new(NestedNode)
	err := finder.Take(&sibling).Error
	if err = IgnoreNotFoundError(err); err != nil {
		return err
	}
	// 重建受影响的左右边界
	if sibling.Depth > 0 {
		n.Lft = sibling.Rgt + 1
	} else if parent != nil {
		n.Lft = parent.Lft + 1
		parent.Rgt += 2 // 上面的数据更新使 parent.Rgt 变成脏数据
	} else {
		n.Lft = 1
	}
	n.Rgt = n.Lft + 1
	if n.Depth > 1 {
		err = MoveEdge(query, n.Lft, "+ 2")
	}
	return err
}

// 左右边界整体移动
func MoveEdge(query *gorm.DB, base uint, offset string) (err error) {
	err = query.Where("rgt >= ?", base).Update("rgt", gorm.Expr("rgt "+offset)).Error
	if err = IgnoreNotFoundError(err); err != nil {
		return
	}
	// 当右边界只有一个变动，说明父类是第一层的最末一个，这时左边界没有符合条件的数据
	if query.RowsAffected > 1 {
		err = query.Where("lft >= ?", base).Update("lft", gorm.Expr("lft "+offset)).Error
		err = IgnoreNotFoundError(err)
	}
	return
}
