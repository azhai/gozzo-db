package export

import (
	"fmt"

	base "github.com/azhai/gozzo-db/construct"
	"github.com/azhai/gozzo-db/schema"
	"github.com/azhai/gozzo-utils/common"
	"github.com/jinzhu/gorm"
)

// 找出Model定义中的注释
func GetTableComments(query *gorm.DB, models ...interface{}) map[string]string {
	result := make(map[string]string)
	for _, value := range models {
		var comment string
		if v, ok := value.(base.ITableComment); ok {
			comment = v.TableComment()
		}
		if comment == "" {
			continue
		}
		scope := query.NewScope(value)
		result[scope.TableName()] = comment
	}
	return result
}

// 补充当前数据库中的表注释，仅当没有时才加上
func FillTableComments(query *gorm.DB, models ...interface{}) *gorm.DB {
	tpl := `
IF NOT EXISTS(
	SELECT NULL FROM INFORMATION_SCHEMA.TABLES
		WHERE table_schema = DATABASE()
		AND table_name = '%s'
		AND table_comment = '')  THEN
	ALTER TABLE %s COMMENT = '%s';
END IF;
`
	comments := GetTableComments(query, models...)
	for name, comment := range comments {
		quoteName := common.WrapWith(name, "`", "`")
		query = query.Exec(fmt.Sprintf(tpl, name, quoteName, comment))
		if err := query.Error; err != nil {
			panic(err)
		}
	}
	return query
}

// 更新当前数据库中的表注释，对比不一致时修改
func AlterTableComments(query *gorm.DB, models ...interface{}) *gorm.DB {
	tpl := `ALTER TABLE %s COMMENT = '%s'`
	comments := GetTableComments(query, models...)
	sch := schema.NewSchema(query.DB())
	tbInfos := sch.ListTable("", false)
	for name, comment := range comments {
		info, ok := tbInfos[name]
		if !ok || info.TableComment == comment {
			continue
		}
		quoteName := sch.Dialect.QuoteIdent(name)
		query = query.Exec(fmt.Sprintf(tpl, quoteName, comment))
		if err := query.Error; err != nil {
			panic(err)
		}
	}
	return query
}
