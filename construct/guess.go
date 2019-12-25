package construct

import (
	"fmt"
	"strings"

	"github.com/azhai/gozzo-db/schema"
	"github.com/azhai/gozzo-utils/common"
)

// 根据数据库字段类型猜测Go类型，适用于mysql或pgsql
func GuessTypeName(ci *schema.ColumnInfo) string {
	switch ci.DatabaseTypeName() {
	default:
		return "string"
	case "BIT", "BINARY", "VARBINARY":
		return "[]byte"
	case "BOOL", "BOOLEAN":
		return "bool"
	case "TINYINT":
		if strings.HasSuffix(ci.FullType, "unsigned") {
			if strings.HasPrefix(ci.FullType, "tinyint(1)") {
				return "bool"
			} else {
				return "uint8"
			}
		} else {
			return "int8"
		}
	case "SMALLINT", "YEAR":
		if strings.HasSuffix(ci.FullType, "unsigned") {
			return "uint16"
		} else {
			return "int16"
		}
	case "INT", "INTEGER", "MEDIUMINT":
		if strings.HasSuffix(ci.FullType, "unsigned") {
			return "uint"
		} else {
			return "int"
		}
	case "BIGINT", "LONGINT", "NUMERIC":
		if strings.HasSuffix(ci.FullType, "unsigned") {
			return "uint64"
		} else {
			return "int64"
		}
	case "DECIMAL", "DOUBLE", "FLOAT", "REAL":
		return "float64"
	case "TIME", "TIMESTAMP", "DATETIME", "DATE":
		return "time.Time"
	}
}

// 根据数据库字段属性设置Struct tags
func GuessStructTags(ci *schema.ColumnInfo) *SqlTag {
	tag := NewSqlTag()
	if ci.IsIndex() {
		if ci.IndexType == schema.PrimaryKey {
			tag.Set("primary_key", "")
		} else if ci.IndexType == schema.UniqueKey {
			tag.Set("unique_index", "")
		} else {
			tag.Set("index", "")
		}
	}

	if ci.IsNotNull() {
		tag.Set("not null", "")
		tag.Set("default", ci.GetDefault())
	}
	if ci.Extra != "" {
		tag.Set(ci.Extra, "")
	}
	if ci.Comment != "" {
		tag.Set("comment", common.WrapWith(ci.Comment, "'", "'"))
	}
	
	// 可能为char、text类型
	size := ci.GetSize()
	if strings.HasPrefix(ci.FullType, "char") && size < 100 {
		tag.Set("type", ci.FullType) 
	} else if strings.HasPrefix(ci.FullType, "enum") {
		tag.Set("type", ci.FullType)
	} else if strings.HasPrefix(ci.FullType, "set") {
		tag.Set("type", ci.FullType)
	} else if strings.Contains(ci.FullType, "text") {
		tag.Set("type", ci.FullType)
	} else if size > 0 && size != 255 { // gorm默认的varchar长度为255，不需要再标注
		tag.Set("size", fmt.Sprintf("%d", size))
	}
	if prec, scale := ci.GetPrecision(); prec > 0 {
		tag.Set("precision", fmt.Sprintf("%d,%d", prec, scale))
	}
	return tag
}
