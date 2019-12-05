package construct

import (
	"fmt"
	"strings"

	"github.com/azhai/gozzo-db/schema"
)

func GuessTypeName(ci *schema.ColumnInfo) string {
	switch ci.DatabaseTypeName() {
	default:
		return "string"
	case "BIT", "BINARY", "VARBINARY":
		return "[]byte"
	case "BOOL", "BOOLEAN":
		return "bool"
	case "TINYINT":
		if strings.HasPrefix(ci.FullType, "tinyint(1)") {
			return "bool"
		} else if strings.HasPrefix(ci.FullType, "tinyint(3)") {
			return "uint8"
		} else {
			return "int"
		}
	case "INT", "INTEGER", "SMALLINT", "YEAR":
		if strings.HasSuffix(ci.FullType, "unsigned") {
			return "uint"
		} else {
			return "int"
		}
	case "BIGINT", "MEDIUMINT", "NUMERIC":
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

func GuessStructTags(ci *schema.ColumnInfo) string {
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
	}
	if ci.Extra != "" {
		tag.Set(ci.Extra, "")
	}
	if ci.Comment != "" {
		tag.Set("comment", ci.Comment)
	}
	// gorm默认的varchar长度为255，不需要再标注
	if size := ci.GetSize(); size > 0 && size != 255 {
		tag.Set("size", fmt.Sprintf("%d", size))
	}
	if size, scale := ci.GetPrecision(); size > 0 {
		tag.Set("precision", fmt.Sprintf("%d,%d", size, scale))
	}
	return tag.String("gorm")
}
