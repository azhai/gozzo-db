package schema

import (
	"database/sql"
	"fmt"
	"strings"
)

type IndexKeyType string

// 索引类型
const (
	IsNotKey   IndexKeyType = ""    // 普通字段
	IndexKey                = "MUL" // 普通索引
	UniqueKey               = "UNI" // 唯一索引
	PrimaryKey              = "PRI" // 主键
)

var DefaultColumnExtra = ColumnExtra{DefaultValue: sql.NullString{}}

type UnknownDriverError struct {
	Driver string
}

func (e UnknownDriverError) Error() string {
	return fmt.Sprintf("unknown database driver %s", e.Driver)
}

// 表名信息
type TableInfo struct {
	DbName       string
	TablePrefix  string
	TableName    string
	TableComment string
	TableRows    int64
	Quote        func(ident string) string
}

func (ti TableInfo) GetFullName(escape bool) string {
	dbname, table := "", ti.TableName
	if escape && ti.Quote != nil {
		table = ti.Quote(ti.TableName)
	}
	if ti.DbName != "" {
		if escape && ti.Quote != nil {
			dbname = ti.Quote(ti.DbName) + "."
		} else {
			dbname = ti.DbName + "."
		}
	}
	return dbname + table
}

type ColumnExtra struct {
	FieldName           string
	FullType            string
	IndexType           IndexKeyType
	DefaultValue        sql.NullString
	MaxSize             int
	PrecSize, PrecScale int
	Extra               string
	Comment             string
}

type ColumnInfo struct {
	Table *TableInfo
	*sql.ColumnType
	ColumnExtra
}

func (ci *ColumnInfo) IsNotNull() bool {
	nullable, ok := ci.Nullable()
	return ok && !nullable
}

func (ci *ColumnInfo) IsPrimaryKey() bool {
	return ci.IndexType == PrimaryKey
}

func (ci *ColumnInfo) IsIndex() bool {
	return ci.IndexType != IsNotKey
}

func (ci *ColumnInfo) GetDefine() string {
	def := ci.FullType
	if def == "" {
		def = ci.GetDbType()
	}
	if ci.IsNotNull() {
		def += " NOT NULL"
	}
	def += fmt.Sprintf(" DEFAULT %s", ci.GetDefault())
	if ci.Extra != "" {
		def += " " + ci.Extra
	}
	return def
}

func (ci *ColumnInfo) GetDefault() string {
	if ci.DefaultValue.Valid {
		return "'" + ci.DefaultValue.String + "'"
	}
	if ci.IsNotNull() {
		return ColumnDefaultValue(ci.DatabaseTypeName())
	}
	return "NULL"
}

func (ci *ColumnInfo) GetSize() int {
	if ci.MaxSize > 0 && ci.MaxSize < 65535 {
		return ci.MaxSize
	}
	if ci.PrecSize > 0 && ci.PrecSize < 65535 {
		return ci.PrecSize
	}
	// go-sql-driver/mysql 不提供数据
	if size, ok := ci.Length(); ok && size < 65535 {
		return int(size)
	}
	return 0
}

func (ci *ColumnInfo) GetPrecision() (int, int) {
	if ci.PrecSize > 0 && ci.PrecSize < 65535 {
		return ci.PrecSize, ci.PrecScale
	}
	// go-sql-driver/mysql 数据错误
	if size, scale, ok := ci.DecimalSize(); ok && size < 65535 {
		return int(size), int(scale)
	}
	return 0, 0
}

func (ci *ColumnInfo) GetDbType() string {
	dbtype := strings.ToLower(ci.DatabaseTypeName())
	if size := ci.GetSize(); size > 0 {
		dbtype += fmt.Sprintf("(%d)", size)
	} else if size, scale := ci.GetPrecision(); size > 0 {
		dbtype += fmt.Sprintf("(%d,%d)", size, scale)
	}
	return dbtype
}

func ColumnDefaultValue(typeName string) string {
	switch typeName {
	default:
		return "''"
	case "BIT", "BINARY", "VARBINARY":
		return "''"
	case "BOOL", "BOOLEAN", "TINYINT":
		return "'0'"
	case "INT", "INTEGER", "SMALLINT", "YEAR":
		return "'0'"
	case "BIGINT", "MEDIUMINT", "NUMERIC":
		return "'0'"
	case "DECIMAL", "DOUBLE", "FLOAT", "REAL":
		return "'0.0'"
	case "TIME", "TIMESTAMP", "DATETIME":
		return "'0000-00-00 00:00:00'"
	case "DATE":
		return "'0000-00-00'"
	}
}
