package cluster

import (
	"fmt"
	"sort"
	"strings"

	base "github.com/azhai/gozzo-db/construct"
	"github.com/azhai/gozzo-db/schema"
	"github.com/jinzhu/gorm"
)

const (
	COUMT_IS_EMPTY   = -1 // 尚未计数
	COUMT_IS_DYNAMIC = -2 // 不缓存计数
)

type ShardingModel interface {
	BaseTableName() string
	TableName() string
}

type Sharding struct {
	tableCounter    map[string]int64
	IsTableNameDesc bool
	DbNameMatch     string
	TableFilter     func(dbname, table string) bool
	*gorm.DB
}

func NewSharding(db *gorm.DB, desc bool) *Sharding {
	return &Sharding{
		tableCounter:    make(map[string]int64),
		IsTableNameDesc: desc,
		DB:              db,
	}
}

// Create insert the value into database
func (s *Sharding) Create(value ShardingModel) *Sharding {
	s.DB = s.DB.Table(value.TableName()).Create(value)
	return s
}

func (s *Sharding) FirstTable(out ShardingModel, isLast bool) (tableName string) {
	baseName := out.BaseTableName()
	tables := s.filterShardingNames(baseName, false)
	if isLast {
		for i := len(tables) - 1; i >= 0; i-- {
			tableName = tables[i]
			if s.getShardingCount(tableName, false) > 0 {
				return
			}
		}
	} else {
		for i := 0; i <= len(tables)-1; i++ {
			tableName = tables[i]
			if s.getShardingCount(tableName, false) > 0 {
				return
			}
		}
	}
	return
}

// 获取符合条件的表名，并统计每张表符合条件的行数
func (s *Sharding) filterShardingNames(baseName string, reload bool) (tables []string) {
	if reload || len(s.tableCounter) == 0 {
		sch := schema.NewSchema(s.DB.DB())
		tbInfos := sch.ListTable(s.DbNameMatch, false)
		for name, info := range tbInfos {
			if baseName != "" && !strings.HasPrefix(name, baseName) {
				continue
			}
			if s.TableFilter != nil && !s.TableFilter(name, info.DbName) {
				continue
			}
			tableName := info.GetFullName(false)
			tables = append(tables, tableName)
			s.tableCounter[tableName] = COUMT_IS_EMPTY
		}
	} else {
		for tableName := range s.tableCounter {
			tables = append(tables, tableName)
		}
	}
	if s.IsTableNameDesc {
		sort.Sort(sort.Reverse(sort.StringSlice(tables)))
	} else {
		sort.Sort(sort.StringSlice(tables))
	}
	return
}

// Count a table
func (s *Sharding) getShardingCount(tableName string, check bool) (count int64) {
	count, ok := s.tableCounter[tableName]
	if check && !ok {
		return
	}
	if count < 0 {
		var oldCount = count
		s.DB.Table(tableName).Count(&count)
		if COUMT_IS_EMPTY == oldCount {
			s.tableCounter[tableName] = count
		}
	}
	return
}

// Count all
func (s *Sharding) CountSharding(model ShardingModel) (count int64) {
	baseName := model.BaseTableName()
	for _, tableName := range s.filterShardingNames(baseName, false) {
		count += s.getShardingCount(tableName, false)
	}
	return
}

// Select step by step.
func (s *Sharding) PaginateSharding(model ShardingModel, page, size int, fetch base.FilterFunc) error {
	if page <= 0 {
		return fmt.Errorf("Param 'page' is out of range")
	}
	if size <= 0 {
		return fmt.Errorf("Param 'size' should be greater than 0")
	}

	remain := size
	offset := int64((page - 1) * size)
	baseName := model.BaseTableName()
	for _, tableName := range s.filterShardingNames(baseName, false) {
		count := s.getShardingCount(tableName, false)
		if offset >= count {
			offset -= count
			continue
		}
		query := s.DB.Table(tableName).Limit(remain)
		if offset > 0 {
			query.Offset(int(offset))
		}
		if err := query.Scopes(fetch).Error; err != nil {
			break
		}
		if remain -= int(count); remain <= 0 {
			break
		}
		offset = 0 // 后续查询不需要偏移了
	}
	return nil
}
